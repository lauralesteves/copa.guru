package repositories

import (
	"context"
	"testing"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func intPtr(v int) *int { return &v }

// setupPredictionRepo returns a repository on a clean test database with the
// unique (userId, matchId) index from docs/mongodb/indexes.mongodb.js.
func setupPredictionRepo(t *testing.T) (*config.MongoContext, PredictionRepository, func()) {
	t.Helper()
	mc, cleanup := setupTestDB(t)

	_, err := mc.Database.Collection(config.PredictionCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "matchId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		cleanup()
		t.Fatalf("create index: %v", err)
	}

	return mc, NewPredictionRepository(mc), cleanup
}

func insertPredictions(t *testing.T, mc *config.MongoContext, predictions ...*models.Prediction) {
	t.Helper()
	docs := make([]any, len(predictions))
	for i, p := range predictions {
		if p.ID.IsZero() {
			p.ID = bson.NewObjectID()
		}
		docs[i] = p
	}
	if _, err := mc.Database.Collection(config.PredictionCollection).InsertMany(context.Background(), docs); err != nil {
		t.Fatalf("insert predictions: %v", err)
	}
}

// --- Upsert ---

func TestPredictionUpsert_CreatesThenUpdates(t *testing.T) {
	_, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	userID, matchID := bson.NewObjectID(), bson.NewObjectID()

	created, err := repo.Upsert(&models.Prediction{UserID: userID, MatchID: matchID, MatchNumber: 7, Goals1: 1, Goals2: 0})
	if err != nil {
		t.Fatalf("Upsert() create error: %v", err)
	}
	if created.ID.IsZero() || created.MatchNumber != 7 || created.Goals1 != 1 || created.Points != nil {
		t.Errorf("created = %+v, want id set, matchNumber 7, goals1 1, no points", created)
	}

	updated, err := repo.Upsert(&models.Prediction{UserID: userID, MatchID: matchID, MatchNumber: 7, Goals1: 2, Goals2: 2})
	if err != nil {
		t.Fatalf("Upsert() update error: %v", err)
	}
	if updated.ID != created.ID {
		t.Errorf("updated.ID = %v, want %v", updated.ID, created.ID)
	}
	if updated.Goals1 != 2 || updated.Goals2 != 2 {
		t.Errorf("updated goals = %d-%d, want 2-2", updated.Goals1, updated.Goals2)
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("createdAt changed on update: %v -> %v", created.CreatedAt, updated.CreatedAt)
	}

	count, err := repo.CountByMatchId(matchID)
	if err != nil {
		t.Fatalf("CountByMatchId() error: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1 prediction per user per match", count)
	}
}

// --- BulkUpsert ---

func TestPredictionBulkUpsert(t *testing.T) {
	_, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	userID := bson.NewObjectID()
	m1, m2 := bson.NewObjectID(), bson.NewObjectID()

	if _, err := repo.Upsert(&models.Prediction{UserID: userID, MatchID: m1, MatchNumber: 1, Goals1: 0, Goals2: 0}); err != nil {
		t.Fatalf("Upsert() error: %v", err)
	}

	saved, err := repo.BulkUpsert([]*models.Prediction{
		{UserID: userID, MatchID: m1, MatchNumber: 1, Goals1: 3, Goals2: 1},
		{UserID: userID, MatchID: m2, MatchNumber: 2, Goals1: 1, Goals2: 1},
	})
	if err != nil {
		t.Fatalf("BulkUpsert() error: %v", err)
	}
	if saved != 2 {
		t.Errorf("saved = %d, want 2", saved)
	}

	p1, _ := repo.GetByUserAndMatch(userID, m1)
	if p1 == nil || p1.Goals1 != 3 || p1.Goals2 != 1 {
		t.Errorf("match 1 prediction = %+v, want 3-1", p1)
	}
	p2, _ := repo.GetByUserAndMatch(userID, m2)
	if p2 == nil || p2.MatchNumber != 2 {
		t.Errorf("match 2 prediction = %+v, want matchNumber 2", p2)
	}
}

func TestPredictionBulkUpsert_Empty(t *testing.T) {
	_, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	saved, err := repo.BulkUpsert(nil)
	if err != nil || saved != 0 {
		t.Errorf("BulkUpsert(nil) = %d, %v, want 0, nil", saved, err)
	}
}

// --- Points ---

func TestPredictionUpdatePoints(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	userID, matchID := bson.NewObjectID(), bson.NewObjectID()
	p := &models.Prediction{UserID: userID, MatchID: matchID, Goals1: 1, Goals2: 0}
	insertPredictions(t, mc, p)

	if err := repo.UpdatePoints(p.ID, 5); err != nil {
		t.Fatalf("UpdatePoints() error: %v", err)
	}

	got, _ := repo.GetByUserAndMatch(userID, matchID)
	if got == nil || got.Points == nil || *got.Points != 5 {
		t.Errorf("points = %v, want 5", got)
	}
}

func TestPredictionBulkUpdatePoints(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	matchID := bson.NewObjectID()
	p1 := &models.Prediction{UserID: bson.NewObjectID(), MatchID: matchID}
	p2 := &models.Prediction{UserID: bson.NewObjectID(), MatchID: matchID}
	insertPredictions(t, mc, p1, p2)

	err := repo.BulkUpdatePoints([]models.PointsUpdate{
		{PredictionID: p1.ID, Points: 10},
		{PredictionID: p2.ID, Points: 0},
	})
	if err != nil {
		t.Fatalf("BulkUpdatePoints() error: %v", err)
	}

	predictions, err := repo.ListByMatchId(matchID)
	if err != nil {
		t.Fatalf("ListByMatchId() error: %v", err)
	}
	points := map[bson.ObjectID]int{}
	for _, p := range predictions {
		if p.Points == nil {
			t.Fatalf("prediction %v has no points", p.ID)
		}
		points[p.ID] = *p.Points
	}
	if points[p1.ID] != 10 || points[p2.ID] != 0 {
		t.Errorf("points = %v, want p1=10 p2=0", points)
	}
}

func TestPredictionBulkUpdatePoints_Empty(t *testing.T) {
	_, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	if err := repo.BulkUpdatePoints(nil); err != nil {
		t.Errorf("BulkUpdatePoints(nil) error: %v", err)
	}
}

// --- Queries ---

func TestPredictionGetByUserAndMatch_NotFound(t *testing.T) {
	_, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	got, err := repo.GetByUserAndMatch(bson.NewObjectID(), bson.NewObjectID())
	if err != nil || got != nil {
		t.Errorf("GetByUserAndMatch() = %v, %v, want nil, nil", got, err)
	}
}

func TestPredictionListByUserId_Filters(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	groupA := &models.Match{ID: bson.NewObjectID(), MatchNumber: 1, Group: "A", Stage: "GROUP_STAGE"}
	groupB := &models.Match{ID: bson.NewObjectID(), MatchNumber: 2, Group: "B", Stage: "GROUP_STAGE"}
	final := &models.Match{ID: bson.NewObjectID(), MatchNumber: 3, Stage: "FINAL"}
	if _, err := mc.Database.Collection(config.MatchCollection).InsertMany(context.Background(), []any{groupA, groupB, final}); err != nil {
		t.Fatalf("insert matches: %v", err)
	}

	userID := bson.NewObjectID()
	insertPredictions(t, mc,
		&models.Prediction{UserID: userID, MatchID: final.ID, MatchNumber: 3},
		&models.Prediction{UserID: userID, MatchID: groupA.ID, MatchNumber: 1},
		&models.Prediction{UserID: userID, MatchID: groupB.ID, MatchNumber: 2},
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: groupA.ID, MatchNumber: 1},
	)

	tests := []struct {
		name    string
		filters models.PredictionFilters
		want    []int
	}{
		{"no filters sorted by match number", models.PredictionFilters{}, []int{1, 2, 3}},
		{"stage", models.PredictionFilters{Stage: "GROUP_STAGE"}, []int{1, 2}},
		{"group", models.PredictionFilters{Group: "B"}, []int{2}},
		{"stage and group", models.PredictionFilters{Stage: "FINAL", Group: "A"}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.ListByUserId(userID, tt.filters)
			if err != nil {
				t.Fatalf("ListByUserId() error: %v", err)
			}
			numbers := make([]int, len(got))
			for i, p := range got {
				numbers[i] = p.MatchNumber
			}
			if len(numbers) != len(tt.want) {
				t.Fatalf("match numbers = %v, want %v", numbers, tt.want)
			}
			for i := range numbers {
				if numbers[i] != tt.want[i] {
					t.Fatalf("match numbers = %v, want %v", numbers, tt.want)
				}
			}
		})
	}
}

// --- Aggregations ---

func TestPredictionGetDistribution(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	matchID := bson.NewObjectID()
	insertPredictions(t, mc,
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: matchID, Goals1: 2, Goals2: 0},
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: matchID, Goals1: 1, Goals2: 0},
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: matchID, Goals1: 1, Goals2: 1},
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: matchID, Goals1: 0, Goals2: 3},
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: bson.NewObjectID(), Goals1: 0, Goals2: 3},
	)

	got, err := repo.GetDistribution(matchID)
	if err != nil {
		t.Fatalf("GetDistribution() error: %v", err)
	}
	want := models.PredictionDistribution{MatchID: matchID.Hex(), Total: 4, HomeWin: 2, Draw: 1, AwayWin: 1}
	if *got != want {
		t.Errorf("GetDistribution() = %+v, want %+v", *got, want)
	}
}

func TestPredictionGetDistribution_Empty(t *testing.T) {
	_, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	matchID := bson.NewObjectID()
	got, err := repo.GetDistribution(matchID)
	if err != nil {
		t.Fatalf("GetDistribution() error: %v", err)
	}
	want := models.PredictionDistribution{MatchID: matchID.Hex()}
	if *got != want {
		t.Errorf("GetDistribution() = %+v, want %+v", *got, want)
	}
}

func TestPredictionGetUserSummary(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	userID := bson.NewObjectID()
	insertPredictions(t, mc,
		&models.Prediction{UserID: userID, MatchID: bson.NewObjectID(), Points: intPtr(10)},
		&models.Prediction{UserID: userID, MatchID: bson.NewObjectID(), Points: intPtr(5)},
		&models.Prediction{UserID: userID, MatchID: bson.NewObjectID(), Points: intPtr(5)},
		&models.Prediction{UserID: userID, MatchID: bson.NewObjectID(), Points: intPtr(0)},
		&models.Prediction{UserID: userID, MatchID: bson.NewObjectID()}, // not scored yet
		&models.Prediction{UserID: bson.NewObjectID(), MatchID: bson.NewObjectID(), Points: intPtr(10)},
	)

	got, err := repo.GetUserSummary(userID)
	if err != nil {
		t.Fatalf("GetUserSummary() error: %v", err)
	}
	want := models.PredictionSummary{Total: 5, ExactScores: 1, CorrectResults: 2, TotalPoints: 20}
	if *got != want {
		t.Errorf("GetUserSummary() = %+v, want %+v", *got, want)
	}
}

func TestPredictionGetUserSummary_Empty(t *testing.T) {
	_, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	got, err := repo.GetUserSummary(bson.NewObjectID())
	if err != nil {
		t.Fatalf("GetUserSummary() error: %v", err)
	}
	if *got != (models.PredictionSummary{}) {
		t.Errorf("GetUserSummary() = %+v, want zero summary", *got)
	}
}

func TestPredictionAggregateLeaderboard(t *testing.T) {
	mc, repo, cleanup := setupPredictionRepo(t)
	defer cleanup()

	alice := &models.User{ID: bson.NewObjectID(), Name: "Alice", Picture: "alice.png"}
	bob := &models.User{ID: bson.NewObjectID(), Name: "Bob"}
	if _, err := mc.Database.Collection(config.UserCollection).InsertMany(context.Background(), []any{alice, bob}); err != nil {
		t.Fatalf("insert users: %v", err)
	}

	// Both have 15 points; Bob wins the tiebreaker with more exact scores.
	insertPredictions(t, mc,
		&models.Prediction{UserID: alice.ID, MatchID: bson.NewObjectID(), Points: intPtr(5)},
		&models.Prediction{UserID: alice.ID, MatchID: bson.NewObjectID(), Points: intPtr(5)},
		&models.Prediction{UserID: alice.ID, MatchID: bson.NewObjectID(), Points: intPtr(5)},
		&models.Prediction{UserID: alice.ID, MatchID: bson.NewObjectID()}, // not scored: ignored
		&models.Prediction{UserID: bob.ID, MatchID: bson.NewObjectID(), Points: intPtr(10)},
		&models.Prediction{UserID: bob.ID, MatchID: bson.NewObjectID(), Points: intPtr(5)},
	)

	entries, err := repo.AggregateLeaderboard()
	if err != nil {
		t.Fatalf("AggregateLeaderboard() error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}

	first, second := entries[0], entries[1]
	if first.UserID != bob.ID || first.UserName != "Bob" || first.TotalPoints != 15 ||
		first.ExactScores != 1 || first.CorrectResults != 1 || first.TotalPredictions != 2 {
		t.Errorf("first = %+v, want Bob with 15 points, 1 exact, 1 correct, 2 predictions", first)
	}
	if second.UserID != alice.ID || second.UserName != "Alice" || second.UserPicture != "alice.png" ||
		second.TotalPoints != 15 || second.ExactScores != 0 || second.CorrectResults != 3 || second.TotalPredictions != 3 {
		t.Errorf("second = %+v, want Alice with 15 points, 0 exact, 3 correct, 3 predictions", second)
	}
}
