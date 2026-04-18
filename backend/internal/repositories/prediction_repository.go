//go:generate mockgen -source=prediction_repository.go -destination=prediction_repository_mock.go -package=repositories

package repositories

import (
	"context"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PredictionRepository interface {
	ListByMatchId(matchID bson.ObjectID) ([]*models.Prediction, error)
	BulkUpdatePoints(updates []models.PointsUpdate) error
	AggregateLeaderboard() ([]*models.LeaderboardEntry, error)

	GetByUserAndMatch(userID, matchID bson.ObjectID) (*models.Prediction, error)
	GetDistribution(matchID bson.ObjectID) (*models.PredictionDistribution, error)
	ListByUserId(userID bson.ObjectID) ([]*models.Prediction, error)
	ListByMatch(matchID bson.ObjectID) ([]*models.Prediction, error)
	BulkUpsert(predictions []*models.Prediction) error
	UpdatePoints(predictionID bson.ObjectID, points int) error
	//GetUserSummary(userID bson.ObjectID) (*models.UserSummary, error)
	CountByMatch(matchID bson.ObjectID) (int, error)
	Upsert(prediction *models.Prediction) (*models.Prediction, error)
}

type predictionRepositoryImpl struct {
	collection *mongo.Collection
}

/* func NewPredictionRepository(db *config.MongoContext) PredictionRepository {
	collection := db.Client.Database(db.Database.Name()).Collection(config.PredictionCollection)

	return &predictionRepositoryImpl{collection: collection}

} */

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

// GetByUserAndMatch gets from Mongo a prediction for a given user and match.
// If no prediction is found, it returns nil and a nil and error.
func (p *predictionRepositoryImpl) GetByUserAndMatch(userID, matchID bson.ObjectID) (*models.Prediction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var prediction models.Prediction
	err := p.collection.FindOne(ctx, bson.M{"userId": userID, "matchId": matchID}).Decode(&prediction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No prediction found for this user and match
		}
		return nil, err
	}

	return &prediction, nil
}
