package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const testDBName = "copa_guru_test"

func setupTestDB(t *testing.T) (*config.MongoContext, func()) {
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

	db := client.Database(testDBName)
	mc := &config.MongoContext{Client: client, Database: db}

	cleanup := func() {
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	}
	return mc, cleanup
}

// --- Get ---

func TestGet(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	now := time.Now()

	user := &models.User{
		Email: "byid@copa.guru",
		Name:  "By ID",
		Auth: models.Auth{
			GoogleID:    "google-byid",
			LastLoginAt: &now,
		},
	}
	created, _ := repo.Upsert(user)

	found, err := repo.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Auth.GoogleID != "google-byid" {
		t.Errorf("GoogleID = %q, want %q", found.Auth.GoogleID, "google-byid")
	}
}

func TestGet_NotFound(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	found, err := repo.Get(bson.NewObjectID())
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent ID")
	}
}

// --- GetByEmail ---

func TestGetByEmail(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	now := time.Now()

	user := &models.User{
		Strategy: models.StrategyEmail,
		Email:    "email@copa.guru",
		Name:     "Email User",
		Auth: models.Auth{
			Password:    "$2a$10$fakehash",
			LastLoginAt: &now,
		},
	}
	repo.Upsert(user)

	found, err := repo.GetByEmail("email@copa.guru")
	if err != nil {
		t.Fatalf("GetByEmail() error: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Strategy != models.StrategyEmail {
		t.Errorf("Strategy = %q, want %q", found.Strategy, models.StrategyEmail)
	}
	if found.Auth.Password != "$2a$10$fakehash" {
		t.Errorf("Password = %q, want %q", found.Auth.Password, "$2a$10$fakehash")
	}
}

func TestGetByEmail_NotFound(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	found, err := repo.GetByEmail("nonexistent@copa.guru")
	if err != nil {
		t.Fatalf("GetByEmail() error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent email")
	}
}

// --- GetByGoogleID ---

func TestGetByGoogleID(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	now := time.Now()

	user := &models.User{
		Email: "find@copa.guru",
		Name:  "Find Me",
		Auth: models.Auth{
			GoogleID:    "google-find",
			LastLoginAt: &now,
		},
	}
	repo.Upsert(user)

	found, err := repo.GetByGoogleID("google-find")
	if err != nil {
		t.Fatalf("GetByGoogleID() error: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Name != "Find Me" {
		t.Errorf("Name = %q, want %q", found.Name, "Find Me")
	}
}

func TestGetByGoogleID_NotFound(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	found, err := repo.GetByGoogleID("nonexistent")
	if err != nil {
		t.Fatalf("GetByGoogleID() error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent user")
	}
}

// --- GetByRefreshToken ---

// (covered indirectly by UpdateRefreshToken + InvalidateRefreshToken tests)

// --- Upsert ---

func TestUpsert_CreateNew(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	now := time.Now()

	user := &models.User{
		Strategy: models.StrategyGoogle,
		Email:    "test@copa.guru",
		Name:     "Test User",
		Picture:  "https://example.com/pic.jpg",
		Auth: models.Auth{
			GoogleID:    "google-123",
			LastLoginAt: &now,
		},
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
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	now := time.Now()

	user := &models.User{
		Strategy: models.StrategyGoogle,
		Email:    "first@copa.guru",
		Name:     "First Name",
		Picture:  "https://example.com/first.jpg",
		Auth: models.Auth{
			GoogleID:    "google-456",
			LastLoginAt: &now,
		},
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

// --- UpdateRefreshToken ---

func TestUpdateRefreshToken(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	now := time.Now()

	user := &models.User{
		Email: "refresh@copa.guru",
		Name:  "Refresh",
		Auth: models.Auth{
			GoogleID:    "google-refresh",
			LastLoginAt: &now,
		},
	}
	created, _ := repo.Upsert(user)

	expires := time.Now().Add(7 * 24 * time.Hour)
	err := repo.UpdateRefreshToken(created.ID, "refresh-token-abc", expires)
	if err != nil {
		t.Fatalf("UpdateRefreshToken() error: %v", err)
	}

	found, _ := repo.Get(created.ID)
	if found.Auth.RefreshToken != "refresh-token-abc" {
		t.Errorf("RefreshToken = %q, want %q", found.Auth.RefreshToken, "refresh-token-abc")
	}
}

// --- InvalidateRefreshToken ---

func TestInvalidateRefreshToken(t *testing.T) {
	mc, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(mc)
	now := time.Now()

	user := &models.User{
		Email: "invalidate@copa.guru",
		Name:  "Invalidate",
		Auth: models.Auth{
			GoogleID:    "google-invalidate",
			LastLoginAt: &now,
		},
	}
	created, _ := repo.Upsert(user)

	expires := time.Now().Add(7 * 24 * time.Hour)
	repo.UpdateRefreshToken(created.ID, "to-be-removed", expires)

	err := repo.InvalidateRefreshToken(created.ID)
	if err != nil {
		t.Fatalf("InvalidateRefreshToken() error: %v", err)
	}

	found, _ := repo.Get(created.ID)
	if found.Auth.RefreshToken != "" {
		t.Errorf("RefreshToken = %q, want empty", found.Auth.RefreshToken)
	}
	if found.Auth.RefreshTokenExpiresAt != nil {
		t.Error("RefreshTokenExpiresAt should be nil")
	}
}
