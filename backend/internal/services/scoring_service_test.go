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

func intPtr(v int) *int { return &v }

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name       string
		predicted1 int
		predicted2 int
		actual1    int
		actual2    int
		want       int
	}{
		{"exact score home win", 2, 1, 2, 1, 10},
		{"exact score away win", 0, 3, 0, 3, 10},
		{"exact score draw", 1, 1, 1, 1, 10},
		{"exact score goalless draw", 0, 0, 0, 0, 10},
		{"correct result home win", 3, 0, 2, 1, 5},
		{"correct result away win", 0, 1, 1, 3, 5},
		{"correct result draw", 1, 1, 0, 0, 5},
		{"correct result high scoring draw", 3, 3, 2, 2, 5},
		{"wrong prediction", 2, 1, 0, 2, 0},
		{"wrong predicted draw actual home win", 1, 1, 2, 0, 0},
		{"wrong predicted home actual away", 2, 0, 0, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePoints(tt.predicted1, tt.predicted2, tt.actual1, tt.actual2)
			if got != tt.want {
				t.Errorf("CalculatePoints(%d, %d, %d, %d) = %d, want %d",
					tt.predicted1, tt.predicted2, tt.actual1, tt.actual2, got, tt.want)
			}
		})
	}
}

func TestScoreMatch(t *testing.T) {
	matchID := bson.NewObjectID()
	matchIDHex := matchID.Hex()

	t.Run("invalid match ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		_, err := svc.ScoreMatch("invalid-id")
		assertServiceError(t, err, svcerr.ErrValidation)
	})

	t.Run("match not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		matchRepo.EXPECT().Get(matchID).Return(nil, nil)

		_, err := svc.ScoreMatch(matchIDHex)
		assertServiceError(t, err, svcerr.ErrNotFound)
	})

	t.Run("match not finished", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		matchRepo.EXPECT().Get(matchID).Return(&models.Match{
			ID:     matchID,
			Status: models.MatchStatusLive,
		}, nil)

		_, err := svc.ScoreMatch(matchIDHex)
		assertServiceError(t, err, svcerr.ErrValidation)
	})

	t.Run("match has no result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		matchRepo.EXPECT().Get(matchID).Return(&models.Match{
			ID:     matchID,
			Status: models.MatchStatusFinished,
		}, nil)

		_, err := svc.ScoreMatch(matchIDHex)
		assertServiceError(t, err, svcerr.ErrValidation)
	})

	t.Run("no predictions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		matchRepo.EXPECT().Get(matchID).Return(&models.Match{
			ID:     matchID,
			Status: models.MatchStatusFinished,
			Goals1: intPtr(2),
			Goals2: intPtr(1),
		}, nil)
		predRepo.EXPECT().ListByMatchId(matchID).Return(nil, nil)

		result, err := svc.ScoreMatch(matchIDHex)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 0 || result.ExactScores != 0 || result.CorrectResults != 0 {
			t.Errorf("expected empty result, got %+v", result)
		}
	})

	t.Run("scores predictions correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		pred1ID := bson.NewObjectID()
		pred2ID := bson.NewObjectID()
		pred3ID := bson.NewObjectID()

		matchRepo.EXPECT().Get(matchID).Return(&models.Match{
			ID:     matchID,
			Status: models.MatchStatusFinished,
			Goals1: intPtr(2),
			Goals2: intPtr(1),
		}, nil)

		predRepo.EXPECT().ListByMatchId(matchID).Return([]*models.Prediction{
			{ID: pred1ID, Goals1: 2, Goals2: 1}, // exact: 10
			{ID: pred2ID, Goals1: 3, Goals2: 0}, // correct result: 5
			{ID: pred3ID, Goals1: 0, Goals2: 2}, // wrong: 0
		}, nil)

		predRepo.EXPECT().BulkUpdatePoints(gomock.Any()).DoAndReturn(func(updates []models.PointsUpdate) error {
			if len(updates) != 3 {
				t.Errorf("expected 3 updates, got %d", len(updates))
			}
			pointsByID := make(map[bson.ObjectID]int)
			for _, u := range updates {
				pointsByID[u.PredictionID] = u.Points
			}
			if pointsByID[pred1ID] != 10 {
				t.Errorf("pred1 expected 10 points, got %d", pointsByID[pred1ID])
			}
			if pointsByID[pred2ID] != 5 {
				t.Errorf("pred2 expected 5 points, got %d", pointsByID[pred2ID])
			}
			if pointsByID[pred3ID] != 0 {
				t.Errorf("pred3 expected 0 points, got %d", pointsByID[pred3ID])
			}
			return nil
		})

		result, err := svc.ScoreMatch(matchIDHex)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 3 {
			t.Errorf("expected total 3, got %d", result.Total)
		}
		if result.ExactScores != 1 {
			t.Errorf("expected 1 exact score, got %d", result.ExactScores)
		}
		if result.CorrectResults != 1 {
			t.Errorf("expected 1 correct result, got %d", result.CorrectResults)
		}
	})

	t.Run("idempotent scoring", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		predID := bson.NewObjectID()
		match := &models.Match{
			ID:     matchID,
			Status: models.MatchStatusFinished,
			Goals1: intPtr(1),
			Goals2: intPtr(1),
		}
		preds := []*models.Prediction{
			{ID: predID, Goals1: 1, Goals2: 1, Points: intPtr(10)},
		}

		matchRepo.EXPECT().Get(matchID).Return(match, nil).Times(2)
		predRepo.EXPECT().ListByMatchId(matchID).Return(preds, nil).Times(2)
		predRepo.EXPECT().BulkUpdatePoints(gomock.Any()).Return(nil).Times(2)

		result1, err := svc.ScoreMatch(matchIDHex)
		if err != nil {
			t.Fatalf("first run: unexpected error: %v", err)
		}

		result2, err := svc.ScoreMatch(matchIDHex)
		if err != nil {
			t.Fatalf("second run: unexpected error: %v", err)
		}

		if result1.Total != result2.Total || result1.ExactScores != result2.ExactScores || result1.CorrectResults != result2.CorrectResults {
			t.Errorf("results differ: run1=%+v run2=%+v", result1, result2)
		}
	})

	t.Run("match repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		matchRepo.EXPECT().Get(matchID).Return(nil, errors.New("db error"))

		_, err := svc.ScoreMatch(matchIDHex)
		assertServiceError(t, err, svcerr.ErrInternal)
	})

	t.Run("prediction repo ListByMatchId error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		matchRepo.EXPECT().Get(matchID).Return(&models.Match{
			ID:     matchID,
			Status: models.MatchStatusFinished,
			Goals1: intPtr(1),
			Goals2: intPtr(0),
		}, nil)
		predRepo.EXPECT().ListByMatchId(matchID).Return(nil, errors.New("db error"))

		_, err := svc.ScoreMatch(matchIDHex)
		assertServiceError(t, err, svcerr.ErrInternal)
	})

	t.Run("bulk update points error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		matchRepo := repositories.NewMockMatchRepository(ctrl)
		predRepo := repositories.NewMockPredictionRepository(ctrl)
		svc := NewScoringService(matchRepo, predRepo)

		matchRepo.EXPECT().Get(matchID).Return(&models.Match{
			ID:     matchID,
			Status: models.MatchStatusFinished,
			Goals1: intPtr(1),
			Goals2: intPtr(0),
		}, nil)
		predRepo.EXPECT().ListByMatchId(matchID).Return([]*models.Prediction{
			{ID: bson.NewObjectID(), Goals1: 1, Goals2: 0},
		}, nil)
		predRepo.EXPECT().BulkUpdatePoints(gomock.Any()).Return(errors.New("db error"))

		_, err := svc.ScoreMatch(matchIDHex)
		assertServiceError(t, err, svcerr.ErrInternal)
	})
}

func assertServiceError(t *testing.T, err error, expectedCode svcerr.ErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var sErr *svcerr.ServiceError
	if !errors.As(err, &sErr) {
		t.Fatalf("expected ServiceError, got %T: %v", err, err)
	}
	if sErr.Code != expectedCode {
		t.Errorf("expected error code %d, got %d: %s", expectedCode, sErr.Code, sErr.Message)
	}
}
