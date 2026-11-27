//go:generate mockgen -source=prediction_service.go -destination=prediction_service_mock.go -package=services
package services

import (
	"time"

	svcerr "github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PredictionService interface {
	BatchUpsert(userID string, predictions []models.PredictionInput) (*models.BatchResult, error)
	ListMyPredictions(userID string, stage string, group string) (*models.MyPredictionsResponse, error)
	ListByMatch(matchNumber int) (*models.MatchPredictionsResponse, error)
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
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, svcerr.NewValidationError("invalid user ID")
	}

	accepted := make([]*models.Prediction, 0, len(predictions))
	rejected := make([]models.RejectedPrediction, 0)

	for idx, input := range predictions {
		if input.Goals1 < 0 || input.Goals1 > 99 || input.Goals2 < 0 || input.Goals2 > 99 {
			rejected = append(rejected, models.RejectedPrediction{
				Index:       idx,
				MatchNumber: input.MatchNumber,
				Reason:      "score needs to be between 0 and 99",
			})
			continue
		}

		match, err := s.matchRepo.FindByMatchNumber(input.MatchNumber)
		if err != nil {
			return nil, svcerr.NewInternalError("failed to find match by number", err)
		}

		if match == nil {
			rejected = append(rejected, models.RejectedPrediction{
				Index:       idx,
				MatchNumber: input.MatchNumber,
				Reason:      "match not found",
			})
			continue
		}

		if !canPredict(match, s.now()) {
			rejected = append(rejected, models.RejectedPrediction{
				Index:       idx,
				MatchNumber: input.MatchNumber,
				Reason:      "prediction for this match already closed", // ou predictions are closed for this match
			})
			continue
		}

		accepted = append(accepted, &models.Prediction{
			UserID:      id,
			MatchID:     match.ID,
			MatchNumber: match.MatchNumber,
			Goals1:      input.Goals1,
			Goals2:      input.Goals2,
		})
	}

	if len(accepted) == 0 {
		return &models.BatchResult{
			Saved:    0,
			Rejected: rejected,
		}, nil
	}

	saved, err := s.predictionRepo.BulkUpsert(accepted)
	if err != nil {
		return nil, svcerr.NewInternalError("failed to save prediction", err)
	}

	return &models.BatchResult{
		Saved:    saved,
		Rejected: rejected,
	}, nil

}

// canPredict check match status and if match still scheduled
func canPredict(match *models.Match, now time.Time) bool {
	return match.Status == models.MatchStatusScheduled && match.Date.After(now)
}

// ListMyPredictions list your predictiosn based on match and userID
func (s *predictionService) ListMyPredictions(userID, stage, group string) (*models.MyPredictionsResponse, error) {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, svcerr.NewValidationError("invalid user ID")
	}

	filters := models.PredictionFilters{
		Stage: stage,
		Group: group,
	}

	records, err := s.predictionRepo.ListByUserWithMatch(id, filters)
	if err != nil {
		return nil, svcerr.NewInternalError("failed to list predictions", err)
	}

	response := &models.MyPredictionsResponse{
		Predictions: make([]*models.PredictionDTO, 0, len(records)),
		Summary: models.PredictionSummary{
			Total: len(records),
		},
	}

	for _, record := range records {
		dto := record.Prediction.ToDTO(record.Match)

		if record.Match == nil || record.Match.Status != models.MatchStatusFinished {
			dto.Points = nil

			if dto.Match != nil {
				dto.Match.Goals1 = nil
				dto.Match.Goals2 = nil
			}

			response.Predictions = append(response.Predictions, dto)
			continue
		}

		if dto.Points != nil {
			points := *dto.Points

			response.Summary.TotalPoints += points

			switch points {
			case models.PointsExactScore:
				response.Summary.ExactScores++
			case models.PointsCorrectResult:
				response.Summary.CorrectResults++
			}

		}

		response.Predictions = append(response.Predictions, dto)
	}

	return response, nil
}

// ListByMatch list match and get users by matchID
func (s *predictionService) ListByMatch(matchNumber int) (*models.MatchPredictionsResponse, error) {
	match, err := s.matchRepo.FindByMatchNumber(matchNumber)
	if err != nil {
		return nil, svcerr.NewInternalError("failed to find match by number", err)
	}
	if match == nil {
		return nil, svcerr.NewNotFoundError("match not found")
	}

	if match.Status != models.MatchStatusFinished {
		return nil, svcerr.NewForbiddenError("predictions are only available after the match has finished")
	}

	records, err := s.predictionRepo.ListByMatchWithUser(match.ID)
	if err != nil {
		return nil, svcerr.NewInternalError("failed to list users prediction in this match", err)
	}

	response := &models.MatchPredictionsResponse{
		Match:       match.ToDTO(),
		Predictions: make([]*models.MatchPredictionDTO, 0, len(records)),
	}

	for _, record := range records {
		item := &models.MatchPredictionDTO{
			UserName:    record.UserName,
			UserPicture: record.UserPicture,
			Prediction:  record.Prediction.ToDTO(nil),
		}
		response.Predictions = append(response.Predictions, item)
	}

	return response, nil
}
