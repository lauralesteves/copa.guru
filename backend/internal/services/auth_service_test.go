package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"github.com/lauralesteves/copa-guru-backend/internal/services/external/google_oauth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"
)

const testSecret = "test-jwt-secret"

func setupAuthService(t *testing.T) (*gomock.Controller, *repositories.MockUserRepository, *google_oauth.MockService, AuthService) {
	ctrl := gomock.NewController(t)
	mockRepo := repositories.NewMockUserRepository(ctrl)
	mockGoogle := google_oauth.NewMockService(ctrl)
	svc := NewAuthService(mockRepo, testSecret, mockGoogle)
	return ctrl, mockRepo, mockGoogle, svc
}

func TestLoginWithGoogle_Success(t *testing.T) {
	ctrl, mockRepo, mockGoogle, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	now := time.Now()

	mockGoogle.EXPECT().Exchange(gomock.Any(), "auth-code", "http://localhost:5173").Return(&google_oauth.GoogleUserInfo{
		GoogleID: "google-123",
		Email:    "test@gmail.com",
		Name:     "Test User",
		Picture:  "https://pic.url",
	}, nil)

	mockRepo.EXPECT().Upsert(gomock.Any()).DoAndReturn(func(user *models.User) (*models.User, error) {
		user.ID = userID
		user.CreatedAt = now
		return user, nil
	})

	mockRepo.EXPECT().UpdateRefreshToken(userID, gomock.Any(), gomock.Any()).Return(nil)

	resp, err := svc.LoginWithGoogle(context.Background(), "auth-code", "http://localhost:5173")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if resp.User.Email != "test@gmail.com" {
		t.Errorf("User.Email = %q, want %q", resp.User.Email, "test@gmail.com")
	}
}

func TestLoginWithGoogle_ExchangeFails(t *testing.T) {
	ctrl, _, mockGoogle, svc := setupAuthService(t)
	defer ctrl.Finish()

	mockGoogle.EXPECT().Exchange(gomock.Any(), "bad-code", "http://localhost").Return(nil, errors.New("invalid code"))

	_, err := svc.LoginWithGoogle(context.Background(), "bad-code", "http://localhost")
	if err == nil {
		t.Error("expected error for failed exchange")
	}
}

func TestLoginWithGoogle_UpsertFails(t *testing.T) {
	ctrl, mockRepo, mockGoogle, svc := setupAuthService(t)
	defer ctrl.Finish()

	mockGoogle.EXPECT().Exchange(gomock.Any(), "code", "http://localhost").Return(&google_oauth.GoogleUserInfo{
		GoogleID: "g-1", Email: "a@b.com", Name: "A",
	}, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil, errors.New("db error"))

	_, err := svc.LoginWithGoogle(context.Background(), "code", "http://localhost")
	if err == nil {
		t.Error("expected error for failed upsert")
	}
}

func TestRefreshTokens_Success(t *testing.T) {
	ctrl, mockRepo, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	expires := time.Now().Add(24 * time.Hour)

	mockRepo.EXPECT().FindByRefreshToken("old-refresh").Return(&models.User{
		ID:                    userID,
		RefreshTokenExpiresAt: &expires,
	}, nil)
	mockRepo.EXPECT().UpdateRefreshToken(userID, gomock.Any(), gomock.Any()).Return(nil)

	resp, err := svc.RefreshTokens("old-refresh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestRefreshTokens_InvalidToken(t *testing.T) {
	ctrl, mockRepo, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByRefreshToken("nonexistent").Return(nil, nil)

	_, err := svc.RefreshTokens("nonexistent")
	if err == nil {
		t.Error("expected error for invalid refresh token")
	}
}

func TestRefreshTokens_Expired(t *testing.T) {
	ctrl, mockRepo, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	expired := time.Now().Add(-1 * time.Hour)

	mockRepo.EXPECT().FindByRefreshToken("expired-token").Return(&models.User{
		ID:                    userID,
		RefreshTokenExpiresAt: &expired,
	}, nil)
	mockRepo.EXPECT().InvalidateRefreshToken(userID).Return(nil)

	_, err := svc.RefreshTokens("expired-token")
	if err == nil {
		t.Error("expected error for expired refresh token")
	}
}

func TestLogout_Success(t *testing.T) {
	ctrl, mockRepo, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	mockRepo.EXPECT().InvalidateRefreshToken(userID).Return(nil)

	err := svc.Logout(userID.Hex())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogout_InvalidID(t *testing.T) {
	ctrl, _, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	err := svc.Logout("not-an-objectid")
	if err == nil {
		t.Error("expected error for invalid user ID")
	}
}

func TestGetMe_Success(t *testing.T) {
	ctrl, mockRepo, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	now := time.Now()

	mockRepo.EXPECT().FindByID(userID).Return(&models.User{
		ID:        userID,
		Email:     "me@copa.guru",
		Name:      "Me",
		CreatedAt: now,
	}, nil)

	resp, err := svc.GetMe(userID.Hex())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Email != "me@copa.guru" {
		t.Errorf("Email = %q, want %q", resp.User.Email, "me@copa.guru")
	}
}

func TestGetMe_NotFound(t *testing.T) {
	ctrl, mockRepo, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	mockRepo.EXPECT().FindByID(userID).Return(nil, nil)

	_, err := svc.GetMe(userID.Hex())
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestGetMe_InvalidID(t *testing.T) {
	ctrl, _, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	_, err := svc.GetMe("bad-id")
	if err == nil {
		t.Error("expected error for invalid user ID")
	}
}

// JWT tests

func TestGenerateAccessToken(t *testing.T) {
	svc := &authService{jwtSecret: testSecret}
	token, err := svc.generateAccessToken("user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestValidateAccessToken_Valid(t *testing.T) {
	svc := &authService{jwtSecret: testSecret}
	token, _ := svc.generateAccessToken("user-123")

	userID, err := svc.validateAccessToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("userID = %q, want %q", userID, "user-123")
	}
}

func TestValidateAccessToken_Expired(t *testing.T) {
	svc := &authService{jwtSecret: testSecret}

	claims := Claims{
		UserID: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))

	_, err := svc.validateAccessToken(signed)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	svc := &authService{jwtSecret: testSecret}
	other := &authService{jwtSecret: "wrong-secret"}

	token, _ := other.generateAccessToken("user-123")
	_, err := svc.validateAccessToken(token)
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestValidateAccessToken_Malformed(t *testing.T) {
	svc := &authService{jwtSecret: testSecret}
	_, err := svc.validateAccessToken("not-a-token")
	if err == nil {
		t.Error("expected error for malformed token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := &authService{jwtSecret: testSecret}
	token1, err := svc.generateRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token1) != 64 {
		t.Errorf("token length = %d, want 64 hex chars", len(token1))
	}

	token2, _ := svc.generateRefreshToken()
	if token1 == token2 {
		t.Error("expected unique tokens")
	}
}
