package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/auth/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const testDBName = "copa_guru_test"
const testCollection = "users"

func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
	t.Helper()
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("MongoDB not reachable: %v", err)
	}

	collection := client.Database(testDBName).Collection(testCollection)

	cleanup := func() {
		_ = client.Database(testDBName).Drop(context.Background())
		_ = client.Disconnect(context.Background())
	}
	return collection, cleanup
}

func TestUpsert_CreateNew(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	now := time.Now()

	user := &models.User{
		Strategy:    models.StrategyGoogle,
		GoogleID:    "google-123",
		Email:       "test@copa.guru",
		Name:        "Test User",
		Picture:     "https://example.com/pic.jpg",
		LastLoginAt: &now,
	}

	result, err := repo.Upsert(user)
	if err != nil {
		t.Fatalf("Upsert() error: %v", err)
	}
	if result.ID.IsZero() {
		t.Error("expected non-zero ID")
	}
	if result.Email != "test@copa.guru" {
		t.Errorf("Email = %q, want %q", result.Email, "test@copa.guru")
	}
}

func TestUpsert_UpdateExisting(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	now := time.Now()

	user := &models.User{
		Strategy:    models.StrategyGoogle,
		GoogleID:    "google-456",
		Email:       "first@copa.guru",
		Name:        "First Name",
		Picture:     "https://example.com/first.jpg",
		LastLoginAt: &now,
	}
	first, _ := repo.Upsert(user)

	user.Name = "Updated Name"
	second, err := repo.Upsert(user)
	if err != nil {
		t.Fatalf("Upsert() error: %v", err)
	}
	if second.ID != first.ID {
		t.Error("expected same ID on upsert update")
	}
	if second.Name != "Updated Name" {
		t.Errorf("Name = %q, want %q", second.Name, "Updated Name")
	}
}

func TestFindByEmail(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	now := time.Now()

	user := &models.User{
		Strategy:    models.StrategyEmail,
		Email:       "email@copa.guru",
		Name:        "Email User",
		Password:    "$2a$10$fakehash",
		LastLoginAt: &now,
	}
	repo.Upsert(user)

	found, err := repo.FindByEmail("email@copa.guru")
	if err != nil {
		t.Fatalf("FindByEmail() error: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Strategy != models.StrategyEmail {
		t.Errorf("Strategy = %q, want %q", found.Strategy, models.StrategyEmail)
	}
	if found.Password != "$2a$10$fakehash" {
		t.Errorf("Password = %q, want %q", found.Password, "$2a$10$fakehash")
	}
}

func TestFindByEmail_NotFound(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	found, err := repo.FindByEmail("nonexistent@copa.guru")
	if err != nil {
		t.Fatalf("FindByEmail() error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent email")
	}
}

func TestFindByGoogleID(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	now := time.Now()

	user := &models.User{
		GoogleID:    "google-find",
		Email:       "find@copa.guru",
		Name:        "Find Me",
		LastLoginAt: &now,
	}
	repo.Upsert(user)

	found, err := repo.FindByGoogleID("google-find")
	if err != nil {
		t.Fatalf("FindByGoogleID() error: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Name != "Find Me" {
		t.Errorf("Name = %q, want %q", found.Name, "Find Me")
	}
}

func TestFindByGoogleID_NotFound(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	found, err := repo.FindByGoogleID("nonexistent")
	if err != nil {
		t.Fatalf("FindByGoogleID() error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent user")
	}
}

func TestFindByID(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	now := time.Now()

	user := &models.User{
		GoogleID:    "google-byid",
		Email:       "byid@copa.guru",
		Name:        "By ID",
		LastLoginAt: &now,
	}
	created, _ := repo.Upsert(user)

	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.GoogleID != "google-byid" {
		t.Errorf("GoogleID = %q, want %q", found.GoogleID, "google-byid")
	}
}

func TestFindByID_NotFound(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	found, err := repo.FindByID(bson.NewObjectID())
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent ID")
	}
}

func TestUpdateRefreshToken(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	now := time.Now()

	user := &models.User{
		GoogleID:    "google-refresh",
		Email:       "refresh@copa.guru",
		Name:        "Refresh",
		LastLoginAt: &now,
	}
	created, _ := repo.Upsert(user)

	expires := time.Now().Add(7 * 24 * time.Hour)
	err := repo.UpdateRefreshToken(created.ID, "refresh-token-abc", expires)
	if err != nil {
		t.Fatalf("UpdateRefreshToken() error: %v", err)
	}

	found, _ := repo.FindByID(created.ID)
	if found.RefreshToken != "refresh-token-abc" {
		t.Errorf("RefreshToken = %q, want %q", found.RefreshToken, "refresh-token-abc")
	}
}

func TestInvalidateRefreshToken(t *testing.T) {
	col, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(col)
	now := time.Now()

	user := &models.User{
		GoogleID:    "google-invalidate",
		Email:       "invalidate@copa.guru",
		Name:        "Invalidate",
		LastLoginAt: &now,
	}
	created, _ := repo.Upsert(user)

	expires := time.Now().Add(7 * 24 * time.Hour)
	repo.UpdateRefreshToken(created.ID, "to-be-removed", expires)

	err := repo.InvalidateRefreshToken(created.ID)
	if err != nil {
		t.Fatalf("InvalidateRefreshToken() error: %v", err)
	}

	found, _ := repo.FindByID(created.ID)
	if found.RefreshToken != "" {
		t.Errorf("RefreshToken = %q, want empty", found.RefreshToken)
	}
	if found.RefreshTokenExpiresAt != nil {
		t.Error("RefreshTokenExpiresAt should be nil")
	}
}
