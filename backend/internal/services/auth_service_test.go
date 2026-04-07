package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"
)

func setupAuthService(t *testing.T) (*gomock.Controller, *repositories.MockUserRepository, *MockJWTService, *MockGoogleOAuthClient, AuthService) {
	ctrl := gomock.NewController(t)
	mockRepo := repositories.NewMockUserRepository(ctrl)
	mockJWT := NewMockJWTService(ctrl)
	mockGoogle := NewMockGoogleOAuthClient(ctrl)
	svc := NewAuthService(mockRepo, mockJWT, mockGoogle)
	return ctrl, mockRepo, mockJWT, mockGoogle, svc
}

func TestLoginWithGoogle_Success(t *testing.T) {
	ctrl, mockRepo, mockJWT, mockGoogle, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	now := time.Now()

	mockGoogle.EXPECT().Exchange(gomock.Any(), "auth-code", "http://localhost:5173").Return(&GoogleUserInfo{
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

	mockJWT.EXPECT().GenerateAccessToken(userID.Hex()).Return("access-token", nil)
	mockJWT.EXPECT().GenerateRefreshToken().Return("refresh-token", nil)
	mockRepo.EXPECT().UpdateRefreshToken(userID, "refresh-token", gomock.Any()).Return(nil)

	resp, err := svc.LoginWithGoogle(context.Background(), "auth-code", "http://localhost:5173")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "access-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "access-token")
	}
	if resp.RefreshToken != "refresh-token" {
		t.Errorf("RefreshToken = %q, want %q", resp.RefreshToken, "refresh-token")
	}
	if resp.User.Email != "test@gmail.com" {
		t.Errorf("User.Email = %q, want %q", resp.User.Email, "test@gmail.com")
	}
}

func TestLoginWithGoogle_ExchangeFails(t *testing.T) {
	ctrl, _, _, mockGoogle, svc := setupAuthService(t)
	defer ctrl.Finish()

	mockGoogle.EXPECT().Exchange(gomock.Any(), "bad-code", "http://localhost").Return(nil, errors.New("invalid code"))

	_, err := svc.LoginWithGoogle(context.Background(), "bad-code", "http://localhost")
	if err == nil {
		t.Error("expected error for failed exchange")
	}
}

func TestLoginWithGoogle_UpsertFails(t *testing.T) {
	ctrl, mockRepo, _, mockGoogle, svc := setupAuthService(t)
	defer ctrl.Finish()

	mockGoogle.EXPECT().Exchange(gomock.Any(), "code", "http://localhost").Return(&GoogleUserInfo{
		GoogleID: "g-1", Email: "a@b.com", Name: "A",
	}, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil, errors.New("db error"))

	_, err := svc.LoginWithGoogle(context.Background(), "code", "http://localhost")
	if err == nil {
		t.Error("expected error for failed upsert")
	}
}

func TestRefreshTokens_Success(t *testing.T) {
	ctrl, mockRepo, mockJWT, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	expires := time.Now().Add(24 * time.Hour)

	mockRepo.EXPECT().FindByRefreshToken("old-refresh").Return(&models.User{
		ID:                    userID,
		RefreshTokenExpiresAt: &expires,
	}, nil)
	mockJWT.EXPECT().GenerateAccessToken(userID.Hex()).Return("new-access", nil)
	mockJWT.EXPECT().GenerateRefreshToken().Return("new-refresh", nil)
	mockRepo.EXPECT().UpdateRefreshToken(userID, "new-refresh", gomock.Any()).Return(nil)

	resp, err := svc.RefreshTokens("old-refresh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "new-access")
	}
	if resp.RefreshToken != "new-refresh" {
		t.Errorf("RefreshToken = %q, want %q", resp.RefreshToken, "new-refresh")
	}
}

func TestRefreshTokens_InvalidToken(t *testing.T) {
	ctrl, mockRepo, _, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByRefreshToken("nonexistent").Return(nil, nil)

	_, err := svc.RefreshTokens("nonexistent")
	if err == nil {
		t.Error("expected error for invalid refresh token")
	}
}

func TestRefreshTokens_Expired(t *testing.T) {
	ctrl, mockRepo, _, _, svc := setupAuthService(t)
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
	ctrl, mockRepo, _, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	mockRepo.EXPECT().InvalidateRefreshToken(userID).Return(nil)

	err := svc.Logout(userID.Hex())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogout_InvalidID(t *testing.T) {
	ctrl, _, _, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	err := svc.Logout("not-an-objectid")
	if err == nil {
		t.Error("expected error for invalid user ID")
	}
}

func TestGetMe_Success(t *testing.T) {
	ctrl, mockRepo, _, _, svc := setupAuthService(t)
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
	ctrl, mockRepo, _, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	userID := bson.NewObjectID()
	mockRepo.EXPECT().FindByID(userID).Return(nil, nil)

	_, err := svc.GetMe(userID.Hex())
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestGetMe_InvalidID(t *testing.T) {
	ctrl, _, _, _, svc := setupAuthService(t)
	defer ctrl.Finish()

	_, err := svc.GetMe("bad-id")
	if err == nil {
		t.Error("expected error for invalid user ID")
	}
}
