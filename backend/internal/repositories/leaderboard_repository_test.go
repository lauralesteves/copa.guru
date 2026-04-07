package repositories

import (
	"testing"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func sampleEntries() []*models.LeaderboardEntry {
	return []*models.LeaderboardEntry{
		{UserID: bson.NewObjectID(), UserName: "Alice", TotalPoints: 30, ExactScores: 2, CorrectResults: 2, TotalPredictions: 5, Rank: 1},
		{UserID: bson.NewObjectID(), UserName: "Bob", TotalPoints: 20, ExactScores: 1, CorrectResults: 2, TotalPredictions: 5, Rank: 2},
		{UserID: bson.NewObjectID(), UserName: "Carol", TotalPoints: 10, ExactScores: 0, CorrectResults: 2, TotalPredictions: 5, Rank: 3},
	}
}

// --- List ---

func TestList(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)
	entries := sampleEntries()
	if err := repo.ReplaceAll(entries); err != nil {
		t.Fatalf("ReplaceAll() error: %v", err)
	}

	results, total, err := repo.List(10, 0)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(results) != 3 {
		t.Errorf("len(results) = %d, want 3", len(results))
	}
	if results[0].Rank != 1 {
		t.Errorf("first result rank = %d, want 1", results[0].Rank)
	}
}

func TestList_Pagination(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)
	if err := repo.ReplaceAll(sampleEntries()); err != nil {
		t.Fatalf("ReplaceAll() error: %v", err)
	}

	results, total, err := repo.List(2, 0)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(results) != 2 {
		t.Errorf("len(results) = %d, want 2", len(results))
	}

	results, _, err = repo.List(2, 2)
	if err != nil {
		t.Fatalf("List() offset error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("len(results) = %d, want 1", len(results))
	}
}

func TestList_Empty(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)

	results, total, err := repo.List(10, 0)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

// --- GetByUserID ---

func TestGetByUserID(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)
	entries := sampleEntries()
	if err := repo.ReplaceAll(entries); err != nil {
		t.Fatalf("ReplaceAll() error: %v", err)
	}

	found, err := repo.GetByUserID(entries[1].UserID)
	if err != nil {
		t.Fatalf("GetByUserID() error: %v", err)
	}
	if found == nil {
		t.Fatal("expected entry, got nil")
	}
	if found.UserName != "Bob" {
		t.Errorf("UserName = %q, want %q", found.UserName, "Bob")
	}
	if found.Rank != 2 {
		t.Errorf("Rank = %d, want 2", found.Rank)
	}
}

func TestGetByUserID_NotFound(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)

	found, err := repo.GetByUserID(bson.NewObjectID())
	if err != nil {
		t.Fatalf("GetByUserID() error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent user")
	}
}

// --- ReplaceAll ---

func TestReplaceAll(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)

	first := sampleEntries()
	if err := repo.ReplaceAll(first); err != nil {
		t.Fatalf("first ReplaceAll() error: %v", err)
	}

	count, _ := repo.Count()
	if count != 3 {
		t.Errorf("count after first replace = %d, want 3", count)
	}

	second := []*models.LeaderboardEntry{
		{UserID: bson.NewObjectID(), UserName: "Dave", TotalPoints: 50, Rank: 1},
		{UserID: bson.NewObjectID(), UserName: "Eve", TotalPoints: 40, Rank: 2},
	}
	if err := repo.ReplaceAll(second); err != nil {
		t.Fatalf("second ReplaceAll() error: %v", err)
	}

	count, _ = repo.Count()
	if count != 2 {
		t.Errorf("count after second replace = %d, want 2", count)
	}

	results, _, _ := repo.List(10, 0)
	if results[0].UserName != "Dave" {
		t.Errorf("first entry = %q, want Dave", results[0].UserName)
	}
}

func TestReplaceAll_Empty(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)

	if err := repo.ReplaceAll(sampleEntries()); err != nil {
		t.Fatalf("ReplaceAll() error: %v", err)
	}

	if err := repo.ReplaceAll([]*models.LeaderboardEntry{}); err != nil {
		t.Fatalf("ReplaceAll(empty) error: %v", err)
	}

	count, _ := repo.Count()
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
}

func TestReplaceAll_SetsUpdatedAt(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)
	entries := []*models.LeaderboardEntry{
		{UserID: bson.NewObjectID(), UserName: "Test", TotalPoints: 10, Rank: 1},
	}

	before := time.Now().Add(-1 * time.Second)
	if err := repo.ReplaceAll(entries); err != nil {
		t.Fatalf("ReplaceAll() error: %v", err)
	}

	results, _, _ := repo.List(1, 0)
	if results[0].UpdatedAt.Before(before) {
		t.Error("expected UpdatedAt to be set")
	}
}

// --- Count ---

func TestCount(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewLeaderboardRepository(mc)

	count, err := repo.Count()
	if err != nil {
		t.Fatalf("Count() error: %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}

	if err := repo.ReplaceAll(sampleEntries()); err != nil {
		t.Fatalf("ReplaceAll() error: %v", err)
	}

	count, err = repo.Count()
	if err != nil {
		t.Fatalf("Count() error: %v", err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}
