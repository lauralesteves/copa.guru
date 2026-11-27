package services

import (
	"errors"
	"reflect"
	"testing"
	"time"

	svcerr "github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"
)

func setupPredictionService(t *testing.T) (*predictionService, *repositories.MockMatchRepository, *repositories.MockPredictionRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	matches := repositories.NewMockMatchRepository(ctrl)
	predictions := repositories.NewMockPredictionRepository(ctrl)
	svc := NewPredictionService(matches, predictions).(*predictionService)
	now := time.Date(2026, 6, 11, 15, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }
	return svc, matches, predictions
}

func TestPredictionBatchUpsert_Lock(t *testing.T) {
	cases := []struct {
		name     string
		status   models.MatchStatus
		delta    time.Duration
		accepted bool
	}{
		{"future scheduled", models.MatchStatusScheduled, time.Hour, true},
		{"kickoff", models.MatchStatusScheduled, 0, false},
		{"past scheduled", models.MatchStatusScheduled, -time.Second, false},
		{"live with future date", models.MatchStatusLive, time.Hour, false},
		{"finished", models.MatchStatusFinished, -time.Hour, false},
		{"postponed", models.MatchStatusPostponed, time.Hour, false},
		{"cancelled", models.MatchStatusCancelled, time.Hour, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, matches, predictions := setupPredictionService(t)
			userID := bson.NewObjectID()
			match := &models.Match{ID: bson.NewObjectID(), MatchNumber: 7, Status: tc.status, Date: svc.now().Add(tc.delta)}
			matches.EXPECT().FindByMatchNumber(7).Return(match, nil)
			if tc.accepted {
				predictions.EXPECT().BulkUpsert(gomock.Any()).DoAndReturn(func(items []*models.Prediction) (int, error) {
					if len(items) != 1 {
						t.Fatalf("saved items = %d, want 1", len(items))
					}
					p := items[0]
					if p.UserID != userID || p.MatchID != match.ID || p.MatchNumber != 7 || p.Goals1 != 2 || p.Goals2 != 1 || p.Points != nil {
						t.Fatalf("unexpected saved prediction: %+v", p)
					}
					return 1, nil
				})
			}
			got, err := svc.BatchUpsert(userID.Hex(), []models.PredictionInput{{MatchNumber: 7, Goals1: 2, Goals2: 1}})
			if err != nil {
				t.Fatal(err)
			}
			if tc.accepted {
				if got.Saved != 1 || len(got.Rejected) != 0 {
					t.Fatalf("result: %+v", got)
				}
			} else {
				if got.Saved != 0 || len(got.Rejected) != 1 {
					t.Fatalf("result: %+v", got)
				}
				if got.Rejected[0].Index != 0 || got.Rejected[0].MatchNumber != 7 || got.Rejected[0].Reason == "" {
					t.Fatalf("rejection: %+v", got.Rejected[0])
				}
			}
		})
	}
}

func TestPredictionBatchUpsert_Goals(t *testing.T) {
	for _, tc := range []struct {
		name  string
		a, b  int
		valid bool
	}{
		{"zero", 0, 0, true}, {"maximum", 99, 99, true}, {"negative home", -1, 0, false}, {"negative away", 0, -1, false}, {"home above maximum", 100, 0, false}, {"away above maximum", 0, 100, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc, matches, predictions := setupPredictionService(t)
			if tc.valid {
				matches.EXPECT().FindByMatchNumber(1).Return(&models.Match{ID: bson.NewObjectID(), MatchNumber: 1, Status: models.MatchStatusScheduled, Date: svc.now().Add(time.Hour)}, nil)
				predictions.EXPECT().BulkUpsert(gomock.Any()).Return(1, nil)
			}
			got, err := svc.BatchUpsert(bson.NewObjectID().Hex(), []models.PredictionInput{{MatchNumber: 1, Goals1: tc.a, Goals2: tc.b}})
			if err != nil {
				t.Fatal(err)
			}
			if tc.valid && (got.Saved != 1 || len(got.Rejected) != 0) {
				t.Fatalf("result: %+v", got)
			}
			if !tc.valid && (got.Saved != 0 || len(got.Rejected) != 1) {
				t.Fatalf("result: %+v", got)
			}
		})
	}
}

func TestPredictionBatchUpsert_Partial(t *testing.T) {
	svc, matches, predictions := setupPredictionService(t)
	future := &models.Match{ID: bson.NewObjectID(), MatchNumber: 1, Status: models.MatchStatusScheduled, Date: svc.now().Add(time.Hour)}
	matches.EXPECT().FindByMatchNumber(1).Return(future, nil)
	matches.EXPECT().FindByMatchNumber(2).Return(&models.Match{Status: models.MatchStatusLive}, nil)
	matches.EXPECT().FindByMatchNumber(3).Return(nil, nil)
	predictions.EXPECT().BulkUpsert(gomock.Any()).DoAndReturn(func(items []*models.Prediction) (int, error) {
		if len(items) != 1 || items[0].MatchID != future.ID {
			t.Fatalf("unexpected accepted items: %+v", items)
		}
		return 1, nil
	})
	got, err := svc.BatchUpsert(bson.NewObjectID().Hex(), []models.PredictionInput{{MatchNumber: 1}, {MatchNumber: 2}, {MatchNumber: 3}, {MatchNumber: 4, Goals1: 100}})
	if err != nil {
		t.Fatal(err)
	}
	if got.Saved != 1 || len(got.Rejected) != 3 {
		t.Fatalf("result: %+v", got)
	}
	for i, r := range got.Rejected {
		if r.Index != i+1 || r.MatchNumber != i+2 || r.Reason == "" {
			t.Errorf("rejection: %+v", r)
		}
	}
}

func TestPredictionBatchUpsert_ErrorsAndEmpty(t *testing.T) {
	t.Run("invalid user", func(t *testing.T) {
		svc, _, _ := setupPredictionService(t)
		got, err := svc.BatchUpsert("bad", nil)
		assertServiceError(t, err, svcerr.ErrValidation)
		if got != nil {
			t.Fatal("expected nil result")
		}
	})
	t.Run("empty batch", func(t *testing.T) {
		svc, _, _ := setupPredictionService(t)
		got, err := svc.BatchUpsert(bson.NewObjectID().Hex(), nil)
		if err != nil || got == nil || got.Saved != 0 || got.Rejected == nil || len(got.Rejected) != 0 {
			t.Fatalf("result=%+v err=%v", got, err)
		}
	})
	for _, writeFailure := range []bool{false, true} {
		name := "match lookup"
		if writeFailure {
			name = "bulk write"
		}
		t.Run(name, func(t *testing.T) {
			svc, matches, predictions := setupPredictionService(t)
			cause := errors.New("database unavailable")
			if writeFailure {
				matches.EXPECT().FindByMatchNumber(1).Return(&models.Match{ID: bson.NewObjectID(), Status: models.MatchStatusScheduled, Date: svc.now().Add(time.Hour)}, nil)
				predictions.EXPECT().BulkUpsert(gomock.Any()).Return(0, cause)
			} else {
				matches.EXPECT().FindByMatchNumber(1).Return(nil, cause)
			}
			got, err := svc.BatchUpsert(bson.NewObjectID().Hex(), []models.PredictionInput{{MatchNumber: 1}})
			assertServiceError(t, err, svcerr.ErrInternal)
			if got != nil || !errors.Is(err, cause) {
				t.Fatalf("result=%+v err=%v", got, err)
			}
		})
	}
}

func TestPredictionListMyPredictions_SummaryAndVisibility(t *testing.T) {
	svc, _, predictions := setupPredictionService(t)
	userID := bson.NewObjectID()
	finished := &models.Match{ID: bson.NewObjectID(), Team1: "BRA", Team2: "JPN", Status: models.MatchStatusFinished, Goals1: intPtr(2), Goals2: intPtr(1)}
	live := &models.Match{Status: models.MatchStatusLive, Goals1: intPtr(1), Goals2: intPtr(0)}
	records := []*models.PredictionWithMatch{
		{Prediction: models.Prediction{ID: bson.NewObjectID(), Points: intPtr(models.PointsExactScore), Goals1: 2, Goals2: 1}, Match: finished},
		{Prediction: models.Prediction{Points: intPtr(models.PointsCorrectResult)}, Match: finished},
		{Prediction: models.Prediction{Points: intPtr(0)}, Match: finished},
		{Prediction: models.Prediction{}, Match: finished},
		{Prediction: models.Prediction{Points: intPtr(10)}, Match: live},
		{Prediction: models.Prediction{Points: intPtr(10)}},
	}
	predictions.EXPECT().ListByUserWithMatch(userID, models.PredictionFilters{Stage: models.StageGroup, Group: "A"}).Return(records, nil)
	got, err := svc.ListMyPredictions(userID.Hex(), models.StageGroup, "A")
	if err != nil {
		t.Fatal(err)
	}
	want := models.PredictionSummary{Total: 6, ExactScores: 1, CorrectResults: 1, TotalPoints: 15}
	if got.Summary != want || len(got.Predictions) != 6 {
		t.Fatalf("result: %+v", got)
	}
	first := got.Predictions[0]
	if first.ID != records[0].ID.Hex() || first.Goals1 != 2 || first.Match == nil || first.Match.Team1 != "BRA" || first.Match.Goals1 == nil || *first.Match.Goals1 != 2 {
		t.Fatalf("first DTO: %+v", first)
	}
	if got.Predictions[2].Points == nil || *got.Predictions[2].Points != 0 {
		t.Fatal("scored zero must remain zero")
	}
	if got.Predictions[3].Points != nil {
		t.Fatal("pending scoring must remain nil")
	}
	if got.Predictions[4].Points != nil || got.Predictions[4].Match.Goals1 != nil || got.Predictions[4].Match.Goals2 != nil {
		t.Fatal("unfinished match exposed result or points")
	}
	if got.Predictions[5].Match != nil || got.Predictions[5].Points != nil {
		t.Fatal("missing match should remain visible without result")
	}
	if records[4].Points == nil || live.Goals1 == nil {
		t.Fatal("response mapping mutated source data")
	}
}

func TestPredictionListMyPredictions_ErrorsAndEmpty(t *testing.T) {
	t.Run("invalid user", func(t *testing.T) {
		svc, _, _ := setupPredictionService(t)
		_, err := svc.ListMyPredictions("bad", "", "")
		assertServiceError(t, err, svcerr.ErrValidation)
	})
	t.Run("empty", func(t *testing.T) {
		svc, _, repo := setupPredictionService(t)
		id := bson.NewObjectID()
		repo.EXPECT().ListByUserWithMatch(id, models.PredictionFilters{}).Return(nil, nil)
		got, err := svc.ListMyPredictions(id.Hex(), "", "")
		if err != nil {
			t.Fatal(err)
		}
		if got.Predictions == nil || len(got.Predictions) != 0 || got.Summary != (models.PredictionSummary{}) {
			t.Fatalf("result: %+v", got)
		}
	})
	t.Run("repository error", func(t *testing.T) {
		svc, _, repo := setupPredictionService(t)
		id := bson.NewObjectID()
		cause := errors.New("db error")
		repo.EXPECT().ListByUserWithMatch(id, models.PredictionFilters{}).Return(nil, cause)
		got, err := svc.ListMyPredictions(id.Hex(), "", "")
		assertServiceError(t, err, svcerr.ErrInternal)
		if got != nil || !errors.Is(err, cause) {
			t.Fatalf("result=%+v err=%v", got, err)
		}
	})
}

func TestPredictionListByMatch_Blocked(t *testing.T) {
	for _, status := range []models.MatchStatus{models.MatchStatusScheduled, models.MatchStatusLive, models.MatchStatusPostponed, models.MatchStatusCancelled} {
		t.Run(string(status), func(t *testing.T) {
			svc, matches, _ := setupPredictionService(t)
			matches.EXPECT().FindByMatchNumber(7).Return(&models.Match{Status: status}, nil)
			got, err := svc.ListByMatch(7)
			assertServiceError(t, err, svcerr.ErrForbidden)
			if got != nil {
				t.Fatal("blocked response must be nil")
			}
			// No prediction expectation: gomock fails if private predictions are queried.
		})
	}
}

func TestPredictionListByMatch_Finished(t *testing.T) {
	for _, empty := range []bool{false, true} {
		name := "with predictions"
		if empty {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			svc, matches, repo := setupPredictionService(t)
			match := &models.Match{ID: bson.NewObjectID(), MatchNumber: 7, Status: models.MatchStatusFinished, Team1: "BRA", Team2: "JPN"}
			records := []*models.PredictionWithUser{}
			if !empty {
				records = append(records, &models.PredictionWithUser{Prediction: models.Prediction{ID: bson.NewObjectID(), Goals1: 2, Goals2: 1, Points: intPtr(5)}, UserName: "Ana", UserPicture: "https://example.com/ana.png"})
			}
			matches.EXPECT().FindByMatchNumber(7).Return(match, nil)
			repo.EXPECT().ListByMatchWithUser(match.ID).Return(records, nil)
			got, err := svc.ListByMatch(7)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.Match, match.ToDTO()) || got.Predictions == nil || len(got.Predictions) != len(records) {
				t.Fatalf("result: %+v", got)
			}
			if !empty {
				item := got.Predictions[0]
				if item.UserName != records[0].UserName || item.UserPicture != records[0].UserPicture || !reflect.DeepEqual(item.Prediction, records[0].Prediction.ToDTO(nil)) {
					t.Fatalf("item: %+v", item)
				}
			}
		})
	}
}

func TestPredictionListByMatch_Errors(t *testing.T) {
	for _, kind := range []string{"lookup", "not found", "predictions"} {
		t.Run(kind, func(t *testing.T) {
			svc, matches, repo := setupPredictionService(t)
			cause := errors.New("db unavailable")
			code := svcerr.ErrInternal
			switch kind {
			case "lookup":
				matches.EXPECT().FindByMatchNumber(7).Return(nil, cause)
			case "not found":
				matches.EXPECT().FindByMatchNumber(7).Return(nil, nil)
				code = svcerr.ErrNotFound
			case "predictions":
				match := &models.Match{ID: bson.NewObjectID(), Status: models.MatchStatusFinished}
				matches.EXPECT().FindByMatchNumber(7).Return(match, nil)
				repo.EXPECT().ListByMatchWithUser(match.ID).Return(nil, cause)
			}
			got, err := svc.ListByMatch(7)
			assertServiceError(t, err, code)
			if got != nil {
				t.Fatal("error response must be nil")
			}
			if kind != "not found" && !errors.Is(err, cause) {
				t.Fatal("original error lost")
			}
		})
	}
}
