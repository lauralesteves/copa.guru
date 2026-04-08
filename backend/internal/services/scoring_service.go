//go:generate mockgen -source=scoring_service.go -destination=scoring_service_mock.go -package=services

package services

import (
	"log/slog"

	svcerr "github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ScoringService interface {
	ScoreMatch(matchID string) (*models.ScoreResult, error)
}

type scoringService struct {
	matchRepo      repositories.MatchRepository
	predictionRepo repositories.PredictionRepository
}

func NewScoringService(matchRepo repositories.MatchRepository, predictionRepo repositories.PredictionRepository) ScoringService {
	return &scoringService{
		matchRepo:      matchRepo,
		predictionRepo: predictionRepo,
	}
}

type matchResult int

const (
	resultHome matchResult = iota
	resultDraw
	resultAway
)

func getResult(goals1, goals2 int) matchResult {
	if goals1 > goals2 {
		return resultHome
	}
	if goals1 < goals2 {
		return resultAway
	}
	return resultDraw
}

func CalculatePoints(predictedGoals1, predictedGoals2, actualGoals1, actualGoals2 int) int {
	if predictedGoals1 == actualGoals1 && predictedGoals2 == actualGoals2 {
		return 10
	}
	if getResult(predictedGoals1, predictedGoals2) == getResult(actualGoals1, actualGoals2) {
		return 5
	}
	return 0
}

func (s *scoringService) ScoreMatch(matchID string) (*models.ScoreResult, error) {
	id, err := bson.ObjectIDFromHex(matchID)
	if err != nil {
		return nil, svcerr.NewValidationError("invalid match ID")
	}

	match, err := s.matchRepo.Get(id)
	if err != nil {
		slog.Error("failed to find match", "matchId", matchID, "error", err)
		return nil, svcerr.NewInternalError("failed to find match", err)
	}
	if match == nil {
		return nil, svcerr.NewNotFoundError("match not found")
	}
	if match.Status != models.MatchStatusFinished {
		return nil, svcerr.NewValidationError("match is not finished")
	}
	if match.Goals1 == nil || match.Goals2 == nil {
		return nil, svcerr.NewValidationError("match has no result")
	}

	predictions, err := s.predictionRepo.ListByMatchId(id)
	if err != nil {
		slog.Error("failed to find predictions", "matchId", matchID, "error", err)
		return nil, svcerr.NewInternalError("failed to find predictions", err)
	}

	result := &models.ScoreResult{}
	var updates []models.PointsUpdate

	for _, p := range predictions {
		points := CalculatePoints(p.Goals1, p.Goals2, *match.Goals1, *match.Goals2)
		updates = append(updates, models.PointsUpdate{
			PredictionID: p.ID,
			Points:       points,
		})
		result.Total++
		if points == 10 {
			result.ExactScores++
		} else if points == 5 {
			result.CorrectResults++
		}
	}

	if len(updates) > 0 {
		if err := s.predictionRepo.BulkUpdatePoints(updates); err != nil {
			slog.Error("failed to bulk update points", "matchId", matchID, "error", err)
			return nil, svcerr.NewInternalError("failed to update prediction points", err)
		}
	}

	slog.Info("scored match", "matchId", matchID, "total", result.Total, "exactScores", result.ExactScores, "correctResults", result.CorrectResults)
	return result, nil
}
