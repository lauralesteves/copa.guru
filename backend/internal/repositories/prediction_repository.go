//go:generate mockgen -source=prediction_repository.go -destination=prediction_repository_mock.go -package=repositories

package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PredictionRepository interface {
	Upsert(prediction *models.Prediction) (*models.Prediction, error)
	BulkUpsert(predictions []*models.Prediction) (int, error)
	UpdatePoints(predictionID bson.ObjectID, points int) error
	BulkUpdatePoints(updates []models.PointsUpdate) error

	GetByUserAndMatch(userID, matchID bson.ObjectID) (*models.Prediction, error)
	ListByUserId(userID bson.ObjectID, filters models.PredictionFilters) ([]*models.Prediction, error)
	ListByMatchId(matchID bson.ObjectID) ([]*models.Prediction, error)
	CountByMatchId(matchID bson.ObjectID) (int64, error)

	GetDistribution(matchID bson.ObjectID) (*models.PredictionDistribution, error)
	GetUserSummary(userID bson.ObjectID) (*models.PredictionSummary, error)
	AggregateLeaderboard() ([]*models.LeaderboardEntry, error)
}

type predictionRepository struct {
	collection *mongo.Collection
}

func NewPredictionRepository(db *config.MongoContext) PredictionRepository {
	collection := db.Client.Database(db.Database.Name()).Collection(config.PredictionCollection)
	return &predictionRepository{collection: collection}
}

// upsertUpdate builds the update that creates or updates a prediction keyed by
// (userId, matchId). Points are never touched: they belong to scoring.
func upsertUpdate(prediction *models.Prediction, now time.Time) bson.M {
	return bson.M{
		"$set": bson.M{
			"goals1":    prediction.Goals1,
			"goals2":    prediction.Goals2,
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"matchNumber": prediction.MatchNumber,
			"createdAt":   now,
		},
	}
}

func userMatchFilter(userID, matchID bson.ObjectID) bson.M {
	return bson.M{"userId": userID, "matchId": matchID}
}

// Upsert creates or updates the prediction of a user for a match, keeping a
// single prediction per (userId, matchId). It returns the stored document.
func (r *predictionRepository) Upsert(prediction *models.Prediction) (*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var saved models.Prediction
	err := r.collection.FindOneAndUpdate(
		ctx,
		userMatchFilter(prediction.UserID, prediction.MatchID),
		upsertUpdate(prediction, time.Now()),
		opts,
	).Decode(&saved)
	if err != nil {
		return nil, err
	}

	return &saved, nil
}

// BulkUpsert creates or updates several predictions in one round trip and
// returns how many were saved (inserted or matched).
func (r *predictionRepository) BulkUpsert(predictions []*models.Prediction) (int, error) {
	if len(predictions) == 0 {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	operations := make([]mongo.WriteModel, 0, len(predictions))
	for _, p := range predictions {
		operations = append(operations, mongo.NewUpdateOneModel().
			SetFilter(userMatchFilter(p.UserID, p.MatchID)).
			SetUpdate(upsertUpdate(p, now)).
			SetUpsert(true))
	}

	result, err := r.collection.BulkWrite(ctx, operations)
	if err != nil {
		return 0, err
	}

	return int(result.UpsertedCount + result.MatchedCount), nil
}

// UpdatePoints sets the points of a single prediction.
func (r *predictionRepository) UpdatePoints(predictionID bson.ObjectID, points int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": predictionID},
		bson.M{"$set": bson.M{"points": points, "updatedAt": time.Now()}},
	)
	return err
}

// BulkUpdatePoints sets the points of several predictions in one round trip.
func (r *predictionRepository) BulkUpdatePoints(updates []models.PointsUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	operations := make([]mongo.WriteModel, 0, len(updates))
	for _, u := range updates {
		operations = append(operations, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": u.PredictionID}).
			SetUpdate(bson.M{"$set": bson.M{"points": u.Points, "updatedAt": now}}))
	}

	_, err := r.collection.BulkWrite(ctx, operations, options.BulkWrite().SetOrdered(false))
	return err
}

// GetByUserAndMatch returns the prediction of a user for a match, or nil when
// there is none.
func (r *predictionRepository) GetByUserAndMatch(userID, matchID bson.ObjectID) (*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var prediction models.Prediction
	err := r.collection.FindOne(ctx, userMatchFilter(userID, matchID)).Decode(&prediction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &prediction, nil
}

// ListByUserId returns the predictions of a user ordered by match number.
// Stage and group filters are applied against the matches collection.
func (r *predictionRepository) ListByUserId(userID bson.ObjectID, filters models.PredictionFilters) ([]*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Aggregate(ctx, listByUserPipeline(userID, filters))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	predictions := []*models.Prediction{}
	if err := cursor.All(ctx, &predictions); err != nil {
		return nil, err
	}

	return predictions, nil
}

func listByUserPipeline(userID bson.ObjectID, filters models.PredictionFilters) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"userId": userID}}},
	}

	if !filters.IsEmpty() {
		matchFilter := bson.M{}
		if filters.Stage != "" {
			matchFilter["match.stage"] = filters.Stage
		}
		if filters.Group != "" {
			matchFilter["match.group"] = filters.Group
		}

		pipeline = append(pipeline,
			bson.D{{Key: "$lookup", Value: bson.M{
				"from":         config.MatchCollection,
				"localField":   "matchId",
				"foreignField": "_id",
				"as":           "match",
			}}},
			bson.D{{Key: "$match", Value: matchFilter}},
			bson.D{{Key: "$project", Value: bson.M{"match": 0}}},
		)
	}

	return append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "matchNumber", Value: 1}}}})
}

// ListByMatchId returns every prediction for a match.
func (r *predictionRepository) ListByMatchId(matchID bson.ObjectID) ([]*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"matchId": matchID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	predictions := []*models.Prediction{}
	if err := cursor.All(ctx, &predictions); err != nil {
		return nil, err
	}

	return predictions, nil
}

// CountByMatchId returns how many predictions a match has.
func (r *predictionRepository) CountByMatchId(matchID bson.ObjectID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{"matchId": matchID})
}

// GetDistribution counts the home win, draw and away win predictions of a
// match. A match without predictions returns all counts as zero.
func (r *predictionRepository) GetDistribution(matchID bson.ObjectID) (*models.PredictionDistribution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	distribution := &models.PredictionDistribution{}
	found, err := aggregateOne(ctx, r.collection, distributionPipeline(matchID), distribution)
	if err != nil {
		return nil, err
	}
	if !found {
		distribution = &models.PredictionDistribution{}
	}

	distribution.MatchID = matchID.Hex()
	return distribution, nil
}

func distributionPipeline(matchID bson.ObjectID) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"matchId": matchID}}},
		{{Key: "$group", Value: bson.M{
			"_id":     nil,
			"total":   bson.M{"$sum": 1},
			"homeWin": countIf(bson.M{"$gt": bson.A{"$goals1", "$goals2"}}),
			"draw":    countIf(bson.M{"$eq": bson.A{"$goals1", "$goals2"}}),
			"awayWin": countIf(bson.M{"$lt": bson.A{"$goals1", "$goals2"}}),
		}}},
	}
}

// GetUserSummary aggregates a user's predictions into totals. A user without
// predictions returns all counts as zero.
func (r *predictionRepository) GetUserSummary(userID bson.ObjectID) (*models.PredictionSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	summary := &models.PredictionSummary{}
	found, err := aggregateOne(ctx, r.collection, userSummaryPipeline(userID), summary)
	if err != nil {
		return nil, err
	}
	if !found {
		return &models.PredictionSummary{}, nil
	}

	return summary, nil
}

func userSummaryPipeline(userID bson.ObjectID) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"userId": userID}}},
		{{Key: "$group", Value: scoreTotals(nil)}},
	}
}

// AggregateLeaderboard sums the scored predictions of every user, joins the
// user name and picture, and sorts by points with exact scores and correct
// results as tiebreakers. Rank is left for the caller to assign.
func (r *predictionRepository) AggregateLeaderboard() ([]*models.LeaderboardEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := r.collection.Aggregate(ctx, leaderboardPipeline())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	entries := []*models.LeaderboardEntry{}
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func leaderboardPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"points": bson.M{"$type": "number"}}}},
		{{Key: "$group", Value: scoreTotals("$userId")}},
		{{Key: "$lookup", Value: bson.M{
			"from":         config.UserCollection,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: "$user"}},
		{{Key: "$project", Value: bson.M{
			"_id":              0,
			"userId":           "$_id",
			"userName":         "$user.name",
			"userPicture":      "$user.picture",
			"totalPoints":      1,
			"exactScores":      1,
			"correctResults":   1,
			"totalPredictions": "$total",
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "totalPoints", Value: -1},
			{Key: "exactScores", Value: -1},
			{Key: "correctResults", Value: -1},
			{Key: "userId", Value: 1},
		}}},
	}
}

// scoreTotals is the $group stage shared by the user summary and the
// leaderboard. $sum ignores predictions without points.
func scoreTotals(groupID any) bson.M {
	return bson.M{
		"_id":            groupID,
		"total":          bson.M{"$sum": 1},
		"totalPoints":    bson.M{"$sum": "$points"},
		"exactScores":    countIf(bson.M{"$eq": bson.A{"$points", models.PointsExactScore}}),
		"correctResults": countIf(bson.M{"$eq": bson.A{"$points", models.PointsCorrectResult}}),
	}
}

func countIf(condition bson.M) bson.M {
	return bson.M{"$sum": bson.M{"$cond": bson.A{condition, 1, 0}}}
}

// aggregateOne runs a pipeline expected to return at most one document and
// decodes it into out. It reports whether a document was found.
func aggregateOne(ctx context.Context, collection *mongo.Collection, pipeline mongo.Pipeline, out any) (bool, error) {
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return false, cursor.Err()
	}

	return true, cursor.Decode(out)
}
