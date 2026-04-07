package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-jwt-secret"

func TestGenerateAccessToken(t *testing.T) {
	svc := NewJWTService(testSecret)
	token, err := svc.GenerateAccessToken("user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestValidateAccessToken_Valid(t *testing.T) {
	svc := NewJWTService(testSecret)
	token, _ := svc.GenerateAccessToken("user-123")

	userID, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("userID = %q, want %q", userID, "user-123")
	}
}

func TestValidateAccessToken_Expired(t *testing.T) {
	svc := &jwtService{secret: testSecret}

	claims := Claims{
		UserID: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))

	_, err := svc.ValidateAccessToken(signed)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	svc := NewJWTService(testSecret)
	other := NewJWTService("wrong-secret")

	token, _ := other.GenerateAccessToken("user-123")
	_, err := svc.ValidateAccessToken(token)
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestValidateAccessToken_Malformed(t *testing.T) {
	svc := NewJWTService(testSecret)
	_, err := svc.ValidateAccessToken("not-a-token")
	if err == nil {
		t.Error("expected error for malformed token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := NewJWTService(testSecret)
	token1, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token1) != 64 {
		t.Errorf("token length = %d, want 64 hex chars", len(token1))
	}

	token2, _ := svc.GenerateRefreshToken()
	if token1 == token2 {
		t.Error("expected unique tokens")
	}
}
