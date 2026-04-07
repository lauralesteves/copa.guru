package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUser_ToDTO(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:      bson.NewObjectID(),
		Email:   "test@copa.guru",
		Name:    "Test User",
		Picture: "https://example.com/pic.jpg",
		Auth: Auth{
			GoogleID:     "google-123",
			RefreshToken: "secret-token",
			LastLoginAt:  &now,
		},
		CreatedAt: now,
	}

	dto := user.ToDTO()

	if dto.ID != user.ID.Hex() {
		t.Errorf("ID = %q, want %q", dto.ID, user.ID.Hex())
	}
	if dto.Email != "test@copa.guru" {
		t.Errorf("Email = %q, want %q", dto.Email, "test@copa.guru")
	}
	if dto.Name != "Test User" {
		t.Errorf("Name = %q, want %q", dto.Name, "Test User")
	}
	if dto.Picture != "https://example.com/pic.jpg" {
		t.Errorf("Picture = %q, want %q", dto.Picture, "https://example.com/pic.jpg")
	}
	if dto.LastLoginAt == nil {
		t.Error("expected LastLoginAt to be set")
	}
}
