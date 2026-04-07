package repositories

import (
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

//go:generate mockgen -source=prediction_repository.go -destination=prediction_repository_mock.go -package=repositories

type PredictionRepository interface {
	FindByMatch(matchID bson.ObjectID) ([]*models.Prediction, error)
	BulkUpdatePoints(updates []models.PointsUpdate) error
}
