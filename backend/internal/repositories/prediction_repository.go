//go:generate mockgen -source=prediction_repository.go -destination=prediction_repository_mock.go -package=repositories

package repositories

import (
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PredictionDistribution struct {
	HomeWin int `json:"homeWin"`
	Draw    int `json:"draw"`
	AwayWin int `json:"awayWin"`
}

type PredictionRepository interface {
	ListByMatchId(matchID bson.ObjectID) ([]*models.Prediction, error)
	BulkUpdatePoints(updates []models.PointsUpdate) error
	AggregateLeaderboard() ([]*models.LeaderboardEntry, error)

	FindByUserAndMatch(userID, matchID bson.ObjectID) (*models.Prediction, error)
	FindByUserId(userID bson.ObjectID) ([]*models.Prediction, error)
	FindByMatch(matchID bson.ObjectID) ([]*models.Prediction, error)
	Upsert(prediction *models.Prediction) (*models.Prediction, error)
	BulkUpsert(predictions []*models.Prediction) error
	UpdatePoints(predictionID bson.ObjectID, points int) error
	//GetUserSummary(userID bson.ObjectID) (*models.UserSummary, error)
	CountByMatch(matchID bson.ObjectID) (int, error)
	GetDistribution(matchID bson.ObjectID) (*PredictionDistribution, error)
}

/* type PredictionRepositoryImpl struct {
	collection *mongo.Collection
} */
