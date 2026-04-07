//go:generate mockgen -source=leaderboard_service.go -destination=leaderboard_service_mock.go -package=services

package services

import (
	"log/slog"

	svcerr "github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type LeaderboardService interface {
	GetLeaderboard(limit, offset int64) ([]*models.LeaderboardEntry, int64, error)
	GetUserPosition(userId string) (*models.LeaderboardEntry, error)
	RebuildLeaderboard() error
}

type leaderboardService struct {
	leaderboardRepo repositories.LeaderboardRepository
	predictionRepo  repositories.PredictionRepository
}

func NewLeaderboardService(leaderboardRepo repositories.LeaderboardRepository, predictionRepo repositories.PredictionRepository) LeaderboardService {
	return &leaderboardService{
		leaderboardRepo: leaderboardRepo,
		predictionRepo:  predictionRepo,
	}
}

// GetLeaderboard returns a paginated list of leaderboard entries sorted by rank.
func (s *leaderboardService) GetLeaderboard(limit, offset int64) ([]*models.LeaderboardEntry, int64, error) {
	entries, total, err := s.leaderboardRepo.List(limit, offset)
	if err != nil {
		slog.Error("failed to list leaderboard", "error", err)
		return nil, 0, svcerr.NewInternalError("failed to get leaderboard", err)
	}
	return entries, total, nil
}

// GetUserPosition returns the leaderboard entry for a specific user.
func (s *leaderboardService) GetUserPosition(userId string) (*models.LeaderboardEntry, error) {
	id, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, svcerr.NewValidationError("invalid user ID")
	}

	entry, err := s.leaderboardRepo.GetByUserID(id)
	if err != nil {
		slog.Error("failed to get user position", "userId", userId, "error", err)
		return nil, svcerr.NewInternalError("failed to get user position", err)
	}
	if entry == nil {
		return nil, svcerr.NewNotFoundError("user not found in leaderboard")
	}

	return entry, nil
}

// RebuildLeaderboard aggregates all scored predictions, ranks users by points
// (tiebreaker: exactScores > correctResults), and replaces the leaderboard collection.
func (s *leaderboardService) RebuildLeaderboard() error {
	entries, err := s.predictionRepo.AggregateLeaderboard()
	if err != nil {
		slog.Error("failed to aggregate leaderboard", "error", err)
		return svcerr.NewInternalError("failed to aggregate leaderboard", err)
	}

	for i, e := range entries {
		e.Rank = i + 1
	}

	if err := s.leaderboardRepo.ReplaceAll(entries); err != nil {
		slog.Error("failed to replace leaderboard", "error", err)
		return svcerr.NewInternalError("failed to replace leaderboard", err)
	}

	slog.Info("leaderboard rebuilt", "totalUsers", len(entries))
	return nil
}
