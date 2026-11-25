package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// setupMatchRepo returns a repository on a clean test database with the
// unique indexes from docs/mongodb/indexes.mongodb.js.
func setupMatchRepo(t *testing.T) (MatchRepository, func()) {
	t.Helper()
	mc, cleanup := setupTestDB(t)

	_, err := mc.Database.Collection(config.MatchCollection).Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "matchNumber", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "externalId", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
	})
	if err != nil {
		cleanup()
		t.Fatalf("create indexes: %v", err)
	}

	return NewMatchRepository(mc), cleanup
}

var kickoff = time.Date(2026, 6, 11, 19, 0, 0, 0, time.UTC)

func newMatch(number int, group, stage string, status models.MatchStatus, date time.Time) *models.Match {
	return &models.Match{
		MatchNumber: number,
		Team1:       "Time A",
		Team2:       "Time B",
		Date:        date,
		Group:       group,
		Stage:       stage,
		Status:      status,
	}
}

func upsertMatches(t *testing.T, repo MatchRepository, matches ...*models.Match) {
	t.Helper()
	for _, m := range matches {
		if err := repo.UpsertByMatchNumber(m); err != nil {
			t.Fatalf("UpsertByMatchNumber(%d) error: %v", m.MatchNumber, err)
		}
	}
}

func matchNumbers(matches []*models.Match) []int {
	numbers := make([]int, len(matches))
	for i, m := range matches {
		numbers[i] = m.MatchNumber
	}
	return numbers
}

func assertNumbers(t *testing.T, got []*models.Match, want ...int) {
	t.Helper()
	numbers := matchNumbers(got)
	if len(numbers) != len(want) {
		t.Fatalf("match numbers = %v, want %v", numbers, want)
	}
	for i := range want {
		if numbers[i] != want[i] {
			t.Fatalf("match numbers = %v, want %v", numbers, want)
		}
	}
}

// --- UpsertByMatchNumber ---

func TestMatchUpsertByMatchNumber_CreatesThenUpdates(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	m := newMatch(1, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff)
	m.ExternalID = "537327"
	m.Stadium = "Estadio Azteca"
	upsertMatches(t, repo, m)

	created, err := repo.FindByMatchNumber(1)
	if err != nil || created == nil {
		t.Fatalf("FindByMatchNumber() = %v, %v", created, err)
	}
	if created.ID.IsZero() || created.Stadium != "Estadio Azteca" || created.ExternalID != "537327" {
		t.Errorf("created = %+v", created)
	}

	m.Team1 = "México"
	m.ExternalID = ""
	upsertMatches(t, repo, m)

	updated, _ := repo.FindByMatchNumber(1)
	if updated.ID != created.ID {
		t.Errorf("ID changed: %s -> %s", created.ID.Hex(), updated.ID.Hex())
	}
	if updated.Team1 != "México" {
		t.Errorf("Team1 = %q, want México", updated.Team1)
	}
	if updated.ExternalID != "537327" {
		t.Errorf("ExternalID = %q, want it kept when upserting without one", updated.ExternalID)
	}
}

func TestMatchUpsertByMatchNumber_KeepsResult(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	upsertMatches(t, repo, newMatch(1, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff))
	if err := repo.UpdateResult(1, 2, 1, models.MatchStatusFinished); err != nil {
		t.Fatalf("UpdateResult() error: %v", err)
	}

	upsertMatches(t, repo, newMatch(1, "A", "GROUP_STAGE", models.MatchStatusFinished, kickoff))

	found, _ := repo.FindByMatchNumber(1)
	if found.Goals1 == nil || found.Goals2 == nil || *found.Goals1 != 2 || *found.Goals2 != 1 {
		t.Errorf("goals = %v-%v, want 2-1 kept", found.Goals1, found.Goals2)
	}
}

func TestMatchUpsertByMatchNumber_WithoutExternalIDs(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	// Two matches without externalId must not clash on the sparse unique index.
	upsertMatches(t, repo,
		newMatch(1, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
		newMatch(2, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
	)

	all, err := repo.FindAll(models.MatchFilters{})
	if err != nil {
		t.Fatalf("FindAll() error: %v", err)
	}
	assertNumbers(t, all, 1, 2)
}

// --- Get / FindByMatchNumber / FindByExternalID ---

func TestMatchGet(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	upsertMatches(t, repo, newMatch(7, "B", "GROUP_STAGE", models.MatchStatusScheduled, kickoff))
	byNumber, _ := repo.FindByMatchNumber(7)

	found, err := repo.Get(byNumber.ID)
	if err != nil || found == nil || found.MatchNumber != 7 {
		t.Fatalf("Get() = %v, %v, want match 7", found, err)
	}
}

func TestMatchFind_NotFound(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	if m, err := repo.Get(bson.NewObjectID()); m != nil || err != nil {
		t.Errorf("Get() = %v, %v, want nil, nil", m, err)
	}
	if m, err := repo.FindByMatchNumber(99); m != nil || err != nil {
		t.Errorf("FindByMatchNumber() = %v, %v, want nil, nil", m, err)
	}
	if m, err := repo.FindByExternalID("missing"); m != nil || err != nil {
		t.Errorf("FindByExternalID() = %v, %v, want nil, nil", m, err)
	}
}

func TestMatchFindByExternalID(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	m := newMatch(3, "C", "GROUP_STAGE", models.MatchStatusScheduled, kickoff)
	m.ExternalID = "537400"
	upsertMatches(t, repo, m, newMatch(4, "C", "GROUP_STAGE", models.MatchStatusScheduled, kickoff))

	found, err := repo.FindByExternalID("537400")
	if err != nil || found == nil || found.MatchNumber != 3 {
		t.Fatalf("FindByExternalID() = %v, %v, want match 3", found, err)
	}

	// An empty id must not match the matches that have none.
	if empty, err := repo.FindByExternalID(""); empty != nil || err != nil {
		t.Errorf("FindByExternalID(\"\") = %v, %v, want nil, nil", empty, err)
	}
}

// --- FindAll ---

func TestMatchFindAll_Filters(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	upsertMatches(t, repo,
		newMatch(3, "B", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
		newMatch(1, "A", "GROUP_STAGE", models.MatchStatusFinished, kickoff),
		newMatch(2, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
		newMatch(104, "", "FINAL", models.MatchStatusScheduled, kickoff),
	)

	tests := []struct {
		name    string
		filters models.MatchFilters
		want    []int
	}{
		{"no filters, ordered by number", models.MatchFilters{}, []int{1, 2, 3, 104}},
		{"group", models.MatchFilters{Group: "A"}, []int{1, 2}},
		{"stage", models.MatchFilters{Stage: "FINAL"}, []int{104}},
		{"status", models.MatchFilters{Status: models.MatchStatusScheduled}, []int{2, 3, 104}},
		{"combined", models.MatchFilters{Group: "A", Status: models.MatchStatusScheduled}, []int{2}},
		{"no match", models.MatchFilters{Group: "Z"}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindAll(tt.filters)
			if err != nil {
				t.Fatalf("FindAll() error: %v", err)
			}
			assertNumbers(t, got, tt.want...)
		})
	}
}

// --- FindScheduledBefore / FindByStatus ---

func TestMatchFindScheduledBefore(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	now := kickoff.Add(3 * time.Hour)
	upsertMatches(t, repo,
		newMatch(1, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff.Add(time.Hour)),
		newMatch(2, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
		newMatch(3, "A", "GROUP_STAGE", models.MatchStatusLive, kickoff),
		newMatch(4, "A", "GROUP_STAGE", models.MatchStatusScheduled, now.Add(time.Hour)),
	)

	got, err := repo.FindScheduledBefore(now)
	if err != nil {
		t.Fatalf("FindScheduledBefore() error: %v", err)
	}
	// Only scheduled matches that already kicked off, ordered by date.
	assertNumbers(t, got, 2, 1)
}

func TestMatchFindByStatus(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	upsertMatches(t, repo,
		newMatch(1, "A", "GROUP_STAGE", models.MatchStatusLive, kickoff.Add(time.Hour)),
		newMatch(2, "A", "GROUP_STAGE", models.MatchStatusFinished, kickoff),
		newMatch(3, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
		newMatch(4, "A", "GROUP_STAGE", models.MatchStatusLive, kickoff),
	)

	got, err := repo.FindByStatus([]models.MatchStatus{models.MatchStatusLive, models.MatchStatusFinished})
	if err != nil {
		t.Fatalf("FindByStatus() error: %v", err)
	}
	assertNumbers(t, got, 2, 4, 1)

	empty, err := repo.FindByStatus(nil)
	if err != nil || empty == nil || len(empty) != 0 {
		t.Errorf("FindByStatus(nil) = %v, %v, want empty slice", empty, err)
	}
}

// --- UpdateResult ---

func TestMatchUpdateResult(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	upsertMatches(t, repo, newMatch(1, "A", "GROUP_STAGE", models.MatchStatusLive, kickoff))

	if err := repo.UpdateResult(1, 0, 0, models.MatchStatusFinished); err != nil {
		t.Fatalf("UpdateResult() error: %v", err)
	}

	found, _ := repo.FindByMatchNumber(1)
	if found.Goals1 == nil || found.Goals2 == nil || *found.Goals1 != 0 || *found.Goals2 != 0 {
		t.Errorf("goals = %v-%v, want 0-0", found.Goals1, found.Goals2)
	}
	if found.Status != models.MatchStatusFinished {
		t.Errorf("status = %s, want FINISHED", found.Status)
	}
}

func TestMatchUpdateResult_NotFound(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	err := repo.UpdateResult(99, 1, 0, models.MatchStatusFinished)
	if !errors.Is(err, ErrMatchNotFound) {
		t.Errorf("UpdateResult() error = %v, want ErrMatchNotFound", err)
	}
}

// --- BulkUpsert ---

func TestMatchBulkUpsert(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	upsertMatches(t, repo, newMatch(1, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff))

	inserted, updated, err := repo.BulkUpsert([]*models.Match{
		newMatch(1, "A", "GROUP_STAGE", models.MatchStatusLive, kickoff),
		newMatch(2, "A", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
		newMatch(3, "B", "GROUP_STAGE", models.MatchStatusScheduled, kickoff),
	})
	if err != nil {
		t.Fatalf("BulkUpsert() error: %v", err)
	}
	if inserted != 2 || updated != 1 {
		t.Errorf("BulkUpsert() = %d inserted, %d updated, want 2, 1", inserted, updated)
	}

	first, _ := repo.FindByMatchNumber(1)
	if first.Status != models.MatchStatusLive {
		t.Errorf("match 1 status = %s, want LIVE", first.Status)
	}
}

func TestMatchBulkUpsert_Empty(t *testing.T) {
	repo, cleanup := setupMatchRepo(t)
	defer cleanup()

	inserted, updated, err := repo.BulkUpsert(nil)
	if err != nil || inserted != 0 || updated != 0 {
		t.Errorf("BulkUpsert(nil) = %d, %d, %v, want 0, 0, nil", inserted, updated, err)
	}
}
