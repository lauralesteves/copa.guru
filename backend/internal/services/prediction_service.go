package services

import (
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
)

type PredictionService interface {
	BatchUpsert(userId string, predictions []models.PredictionInput) (*models.BatchResult, error)
	//ListMyPredictions(userId string, stage, group) (*MyPredictionsResponse, error)
	//ListByMatch(matchNumber) (*MatchPredictionsResponse, error)
}

type predictionService struct {
	matchRepo      repositories.MatchRepository
	predictionRepo repositories.PredictionRepository
	now            func() time.Time
}

func NewPredictionService(
	matchRepo repositories.MatchRepository,
	predictionRepo repositories.PredictionRepository,
) PredictionService {
	return &predictionService{
		matchRepo:      matchRepo,
		predictionRepo: predictionRepo,
		now:            time.Now,
	}
}

func (s *predictionService) BatchUpsert(userID string, predictions []models.PredictionInput) (*models.BatchResult, error) {
	panic("BatchUpsert vai ser feito, para de dar erro caraio")
}
