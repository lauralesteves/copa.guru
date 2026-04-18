//go:generate mockgen -source=prediction_repository.go -destination=prediction_repository_mock.go -package=repositories

package repositories

import (
	"context"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PredictionRepository interface {
	Upsert(prediction *models.Prediction) (*models.Prediction, error)
	BulkUpsert(predictions []*models.Prediction) error
	UpdatePoints(predictionID bson.ObjectID, points int) error
	BulkUpdatePoints(updates []models.PointsUpdate) error

	GetByUserAndMatch(userID, matchID bson.ObjectID) (*models.Prediction, error)
	ListByUserId(userID bson.ObjectID) ([]*models.Prediction, error)
	ListByMatchId(matchID bson.ObjectID) ([]*models.Prediction, error)
	CountByMatchId(matchID bson.ObjectID) (int, error)

	GetDistribution(matchID bson.ObjectID) (*models.PredictionDistribution, error)
	GetUserSummary(userID bson.ObjectID) (*models.UserSummary, error)
	AggregateLeaderboard() ([]*models.LeaderboardEntry, error)
}

type predictionRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPredictionRepository(db *config.MongoContext) PredictionRepository {
	collection := db.Client.Database(db.Database.Name()).Collection(config.PredictionCollection)

	return &predictionRepositoryImpl{collection: collection}

}

// Upsert creates or updates a prediction for a given user and match. It returns the updated prediction.
func (p *predictionRepositoryImpl) Upsert(prediction *models.Prediction) (*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"userId":  prediction.UserID,
		"matchId": prediction.MatchID,
	}

	t := time.Now()

	update := bson.M{
		"$set": bson.M{
			"goals1":    prediction.Goals1,
			"goals2":    prediction.Goals2,
			"updatedAt": t,
		},
		"$setOnInsert": bson.M{
			"matchNumber": prediction.MatchNumber,
			"createdAt":   t,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var updatedPrediction models.Prediction
	err := p.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPrediction)
	if err != nil {
		return nil, err
	}

	return &updatedPrediction, nil
}

// BulkUpsert creates or updates multiple predictions for given users and matches.
// It returns an error if any operation fails.
func (p *predictionRepositoryImpl) BulkUpsert(predictions []*models.Prediction) error {
	if len(predictions) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	operations := make([]mongo.WriteModel, len(predictions))
	for _, prediction := range predictions {
		filter := bson.M{
			"userId":  prediction.UserID,
			"matchId": prediction.MatchID,
		}

		update := bson.M{
			"$set": bson.M{
				"goals1":    prediction.Goals1,
				"goals2":    prediction.Goals2,
				"updatedAt": time.Now(),
			},
			"$setOnInsert": bson.M{
				"createdAt": time.Now(),
			},
		}

		op := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		operations = append(operations, op)
	}

	_, err := p.collection.BulkWrite(ctx, operations)
	return err
}

// UpdatePoints updates the points for a given prediction ID.
func (p *predictionRepositoryImpl) UpdatePoints(predictionID bson.ObjectID, points int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := p.collection.UpdateOne(ctx, bson.M{"_id": predictionID}, bson.M{"$set": bson.M{"points": points, "updatedAt": time.Now()}})
	return err
}

// BulkUpdatePoints updates the points for multiple predictions in a single operation.
func (p *predictionRepositoryImpl) BulkUpdatePoints(updates []models.PointsUpdate) error {

	if len(updates) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	operations := make([]mongo.WriteModel, len(updates))
	for _, update := range updates {
		op := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": update.PredictionID}).
			SetUpdate(bson.M{"$set": bson.M{"points": update.Points, "updatedAt": time.Now()}})
		operations = append(operations, op)
	}

	_, err := p.collection.BulkWrite(ctx, operations)
	return err
}

// GetByUserAndMatch gets from Mongo a prediction for a given user and match.
// If no prediction is found, it returns nil and a nil and error.
func (p *predictionRepositoryImpl) GetByUserAndMatch(userID, matchID bson.ObjectID) (*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var prediction models.Prediction
	err := p.collection.FindOne(ctx, bson.M{"userId": userID, "matchId": matchID}).Decode(&prediction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &prediction, nil
}

// ListByUserId returns all predictions for a given user ID, returns the first created to the oldest.
func (p *predictionRepositoryImpl) ListByUserId(userID bson.ObjectID) ([]*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := p.collection.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*models.Prediction
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// ListByMatchId returns all predictions for a given match ID.
// If no predictions are found, it returns an empty slice and a nil error.
func (p *predictionRepositoryImpl) ListByMatchId(matchID bson.ObjectID) ([]*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := p.collection.Find(ctx, bson.M{"matchId": matchID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*models.Prediction
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// CountByMatchId return the count in int of matches by ID.
func (p *predictionRepositoryImpl) CountByMatchId(matchID bson.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	total, err := p.collection.CountDocuments(ctx, bson.M{"matchId": matchID})
	if err != nil {
		return 0, err
	}

	return int(total), nil
}

// GetDistribution returns the distribution of predictions for a given match ID, including the count of home wins, draws, and away wins.
// If no predictions are found, it returns a distribution with all counts set to zero.
func (p *predictionRepositoryImpl) GetDistribution(matchID bson.ObjectID) (*models.PredictionDistribution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"matchId": matchID}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,

			"homeWin": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$gt": bson.A{"$goals1", "$goals2"}},
						1,
						0,
					},
				},
			},

			"draw": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$goals1", "$goals2"}},
						1,
						0,
					},
				},
			},

			"awayWin": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$lt": bson.A{"$goals1", "$goals2"}},
						1,
						0,
					},
				},
			},
		}}},
	}

	cursor, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*models.PredictionDistribution
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.PredictionDistribution{
			HomeWin: 0,
			Draw:    0,
			AwayWin: 0,
		}, nil
	}

	return results[0], nil
}

// GetUserSummary returns a summary of a user's predictions, including total points, exact scores, and correct results.
func (p *predictionRepositoryImpl) GetUserSummary(userID bson.ObjectID) (*models.UserSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"userId": userID}}},
		{{Key: "$group", Value: bson.M{
			"_id":            "$userId",
			"totalPoints":    bson.M{"$sum": "$points"},
			"exactScores":    bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$points", 3}}, 1, 0}}},
			"correctResults": bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$gt": bson.A{"$points", 0}}, 1, 0}}},
		}}},
	}

	cursor, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []*models.UserSummary
	if err := cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}

	if len(summaries) == 0 {
		return &models.UserSummary{
			UserID:         userID.Hex(),
			TotalPoints:    0,
			ExactScores:    0,
			CorrectResults: 0,
		}, nil
	}

	return summaries[0], nil
}

// AggregateLeaderboard check the total points for each user and return a sorted -1 list of users and their total points.
func (p *predictionRepositoryImpl) AggregateLeaderboard() ([]*models.LeaderboardEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":    "$userId",
			"points": bson.M{"$sum": "$points"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: "$user"}},
		{{Key: "$project", Value: bson.M{
			"userId":   "$_id",
			"username": "$user.username",
			"points":   1,
			"_id":      0,
		}}},
		{{Key: "$sort", Value: bson.M{"points": -1}}},
	}

	cursor, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leaderboard []*models.LeaderboardEntry
	if err := cursor.All(ctx, &leaderboard); err != nil {
		return nil, err
	}

	return leaderboard, nil
}
