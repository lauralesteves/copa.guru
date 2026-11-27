package repositories

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// These integration tests exercise the actual lookup and BSON decoding.
// Like the existing repository tests, they skip when MongoDB is unavailable.
func TestPredictionListByUserWithMatch(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()
	userID := bson.NewObjectID()
	matchA := &models.Match{ID: bson.NewObjectID(), MatchNumber: 1, Team1: "BRA", Team2: "JPN", Group: "A", Stage: models.StageGroup, Status: models.MatchStatusFinished, Goals1: intPtr(2), Goals2: intPtr(1), Date: time.Date(2026, 6, 11, 15, 0, 0, 0, time.UTC)}
	matchB := &models.Match{ID: bson.NewObjectID(), MatchNumber: 2, Group: "B", Stage: models.StageGroup}
	final := &models.Match{ID: bson.NewObjectID(), MatchNumber: 3, Stage: models.StageFinal}
	if _, err := mc.Database.Collection(config.MatchCollection).InsertMany(context.Background(), []any{matchA, matchB, final}); err != nil {
		t.Fatal(err)
	}
	insertPredictions(t, mc,
		&models.Prediction{UserID: userID, MatchID: final.ID, MatchNumber: 3},
		&models.Prediction{UserID: userID, MatchID: matchB.ID, MatchNumber: 2},
		&models.Prediction{UserID: userID, MatchID: matchA.ID, MatchNumber: 1, Goals1: 2, Goals2: 1, Points: intPtr(10)},
		&models.Prediction{UserID: userID, MatchID: bson.NewObjectID(), MatchNumber: 4},
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: matchA.ID, MatchNumber: 1},
	)
	for _, tc := range []struct {
		name    string
		filters models.PredictionFilters
		want    []int
	}{
		{"unfiltered includes missing match", models.PredictionFilters{}, []int{1, 2, 3, 4}},
		{"stage", models.PredictionFilters{Stage: models.StageGroup}, []int{1, 2}},
		{"group", models.PredictionFilters{Group: "A"}, []int{1}},
		{"both", models.PredictionFilters{Stage: models.StageGroup, Group: "A"}, []int{1}},
		{"no matches", models.PredictionFilters{Stage: models.StageFinal, Group: "A"}, []int{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.ListByUserWithMatch(userID, tc.filters)
			if err != nil {
				t.Fatal(err)
			}
			numbers := make([]int, 0, len(got))
			for _, record := range got {
				numbers = append(numbers, record.MatchNumber)
				if record.UserID != userID {
					t.Fatal("another user's prediction leaked")
				}
				if record.MatchNumber == 1 {
					if record.Match == nil || record.Match.ID != matchA.ID || record.Match.Team1 != "BRA" || record.Match.Goals1 == nil || *record.Match.Goals1 != 2 || !record.Match.Date.Equal(matchA.Date) || record.Points == nil || *record.Points != 10 || record.Goals1 != 2 {
						t.Fatalf("incorrect enriched record: %+v", record)
					}
				}
				if record.MatchNumber == 4 && record.Match != nil {
					t.Fatal("missing match should decode as nil")
				}
			}
			if !reflect.DeepEqual(numbers, tc.want) {
				t.Fatalf("numbers=%v want=%v", numbers, tc.want)
			}
		})
	}
	got, err := repo.ListByUserWithMatch(bson.NewObjectID(), models.PredictionFilters{})
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("empty result=%+v err=%v", got, err)
	}
}

func TestPredictionListByMatchWithUser(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()
	matchID := bson.NewObjectID()
	user := &models.User{ID: bson.NewObjectID(), Name: "Ana", Picture: "https://example.com/ana.png", Email: "private@example.com"}
	if _, err := mc.Database.Collection(config.UserCollection).InsertOne(context.Background(), user); err != nil {
		t.Fatal(err)
	}
	firstID := bson.ObjectID{11: 1}
	secondID := bson.ObjectID{11: 2}
	insertPredictions(t, mc,
		&models.Prediction{ID: secondID, UserID: bson.NewObjectID(), MatchID: matchID, MatchNumber: 7, Goals1: 1},
		&models.Prediction{ID: firstID, UserID: user.ID, MatchID: matchID, MatchNumber: 7, Goals1: 2, Goals2: 1, Points: intPtr(5)},
		&models.Prediction{UserID: user.ID, MatchID: bson.NewObjectID(), MatchNumber: 8},
	)
	got, err := repo.ListByMatchWithUser(matchID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("records=%d want=2", len(got))
	}
	if got[0].ID != firstID || got[1].ID != secondID {
		t.Fatal("results not sorted by prediction ID")
	}
	if got[0].UserName != user.Name || got[0].UserPicture != user.Picture || got[0].MatchID != matchID || got[0].Goals1 != 2 || got[0].Goals2 != 1 || got[0].Points == nil || *got[0].Points != 5 {
		t.Fatalf("incorrect enriched record: %+v", got[0])
	}
	if got[1].UserName != "" || got[1].UserPicture != "" {
		t.Fatalf("missing user should have empty public fields: %+v", got[1])
	}
	// Inspect the actual projection result, including fields ignored by the Go model.
	cursor, err := mc.Database.Collection(config.PredictionCollection).Aggregate(context.Background(), listByMatchWithUserPipeline(matchID))
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(context.Background())
	var docs []bson.M
	if err := cursor.All(context.Background(), &docs); err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{"_id": true, "userId": true, "matchId": true, "matchNumber": true, "goals1": true, "goals2": true, "points": true, "createdAt": true, "updatedAt": true, "userName": true, "userPicture": true}
	for _, doc := range docs {
		for key := range doc {
			if !allowed[key] {
				t.Errorf("unexpected projected field %q", key)
			}
		}
	}
	empty, err := repo.ListByMatchWithUser(bson.NewObjectID())
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("empty result=%+v err=%v", empty, err)
	}
}
