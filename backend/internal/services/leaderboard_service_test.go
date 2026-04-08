package services

import (
	"errors"
	"testing"

	svcerr "github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"
)

func setupLeaderboardService(t *testing.T) (*gomock.Controller, *repositories.MockLeaderboardRepository, *repositories.MockPredictionRepository, LeaderboardService) {
	ctrl := gomock.NewController(t)
	mockLeaderboard := repositories.NewMockLeaderboardRepository(ctrl)
	mockPrediction := repositories.NewMockPredictionRepository(ctrl)
	svc := NewLeaderboardService(mockLeaderboard, mockPrediction)
	return ctrl, mockLeaderboard, mockPrediction, svc
}

// --- GetLeaderboard ---

func TestGetLeaderboard_Success(t *testing.T) {
	ctrl, mockLeaderboard, _, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	entries := []*models.LeaderboardEntry{
		{UserName: "Alice", TotalPoints: 30, Rank: 1},
		{UserName: "Bob", TotalPoints: 20, Rank: 2},
	}

	mockLeaderboard.EXPECT().List(int64(10), int64(0)).Return(entries, int64(2), nil)

	result, total, err := svc.GetLeaderboard(10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}
}

func TestGetLeaderboard_RepoError(t *testing.T) {
	ctrl, mockLeaderboard, _, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	mockLeaderboard.EXPECT().List(int64(10), int64(0)).Return(nil, int64(0), errors.New("db error"))

	_, _, err := svc.GetLeaderboard(10, 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- GetUserPosition ---

func TestGetUserPosition_Success(t *testing.T) {
	ctrl, mockLeaderboard, _, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	userId := bson.NewObjectID()
	entry := &models.LeaderboardEntry{UserID: userId, UserName: "Alice", Rank: 1, TotalPoints: 30}

	mockLeaderboard.EXPECT().GetByUserID(userId).Return(entry, nil)

	result, err := svc.GetUserPosition(userId.Hex())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UserName != "Alice" {
		t.Errorf("UserName = %q, want Alice", result.UserName)
	}
}

func TestGetUserPosition_InvalidID(t *testing.T) {
	ctrl, _, _, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	_, err := svc.GetUserPosition("bad-id")
	assertServiceError(t, err, svcerr.ErrValidation)
}

func TestGetUserPosition_NotFound(t *testing.T) {
	ctrl, mockLeaderboard, _, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	userId := bson.NewObjectID()
	mockLeaderboard.EXPECT().GetByUserID(userId).Return(nil, nil)

	_, err := svc.GetUserPosition(userId.Hex())
	assertServiceError(t, err, svcerr.ErrNotFound)
}

func TestGetUserPosition_RepoError(t *testing.T) {
	ctrl, mockLeaderboard, _, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	userId := bson.NewObjectID()
	mockLeaderboard.EXPECT().GetByUserID(userId).Return(nil, errors.New("db error"))

	_, err := svc.GetUserPosition(userId.Hex())
	assertServiceError(t, err, svcerr.ErrInternal)
}

// --- RebuildLeaderboard ---

func TestRebuildLeaderboard_Success(t *testing.T) {
	ctrl, mockLeaderboard, mockPrediction, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	entries := []*models.LeaderboardEntry{
		{UserName: "Alice", TotalPoints: 30},
		{UserName: "Bob", TotalPoints: 20},
	}

	mockPrediction.EXPECT().AggregateLeaderboard().Return(entries, nil)
	mockLeaderboard.EXPECT().ReplaceAll(gomock.Any()).DoAndReturn(func(e []*models.LeaderboardEntry) error {
		if len(e) != 2 {
			t.Errorf("expected 2 entries, got %d", len(e))
		}
		if e[0].Rank != 1 {
			t.Errorf("first rank = %d, want 1", e[0].Rank)
		}
		if e[1].Rank != 2 {
			t.Errorf("second rank = %d, want 2", e[1].Rank)
		}
		return nil
	})

	err := svc.RebuildLeaderboard()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRebuildLeaderboard_Empty(t *testing.T) {
	ctrl, mockLeaderboard, mockPrediction, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	mockPrediction.EXPECT().AggregateLeaderboard().Return([]*models.LeaderboardEntry{}, nil)
	mockLeaderboard.EXPECT().ReplaceAll(gomock.Any()).Return(nil)

	err := svc.RebuildLeaderboard()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRebuildLeaderboard_AggregateError(t *testing.T) {
	ctrl, _, mockPrediction, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	mockPrediction.EXPECT().AggregateLeaderboard().Return(nil, errors.New("db error"))

	err := svc.RebuildLeaderboard()
	assertServiceError(t, err, svcerr.ErrInternal)
}

func TestRebuildLeaderboard_ReplaceError(t *testing.T) {
	ctrl, mockLeaderboard, mockPrediction, svc := setupLeaderboardService(t)
	defer ctrl.Finish()

	mockPrediction.EXPECT().AggregateLeaderboard().Return([]*models.LeaderboardEntry{{UserName: "A"}}, nil)
	mockLeaderboard.EXPECT().ReplaceAll(gomock.Any()).Return(errors.New("db error"))

	err := svc.RebuildLeaderboard()
	assertServiceError(t, err, svcerr.ErrInternal)
}
