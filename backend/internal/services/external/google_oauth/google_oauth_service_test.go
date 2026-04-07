package google_oauth

import (
	"errors"
	"testing"

	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.uber.org/mock/gomock"
)

func TestExchange_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAdapter := NewMockAdapter(ctrl)
	svc := NewService(mockAdapter, "client-id", "client-secret")

	dto := &models.GoogleLoginRequestDTO{Code: "auth-code", RedirectURI: "http://localhost:5173"}

	mockAdapter.EXPECT().ExchangeToken(gomock.Any(), "auth-code", "client-id", "client-secret", "http://localhost:5173").
		Return(&TokenResponse{AccessToken: "google-access-token"}, nil)

	mockAdapter.EXPECT().GetUserInfo(gomock.Any(), "google-access-token").
		Return(&UserInfoResponse{
			Sub:     "google-123",
			Email:   "test@gmail.com",
			Name:    "Test User",
			Picture: "https://pic.url",
		}, nil)

	info, err := svc.Exchange(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.GoogleID != "google-123" {
		t.Errorf("GoogleID = %q, want %q", info.GoogleID, "google-123")
	}
	if info.Email != "test@gmail.com" {
		t.Errorf("Email = %q, want %q", info.Email, "test@gmail.com")
	}
	if info.Name != "Test User" {
		t.Errorf("Name = %q, want %q", info.Name, "Test User")
	}
	if info.Picture != "https://pic.url" {
		t.Errorf("Picture = %q, want %q", info.Picture, "https://pic.url")
	}
}

func TestExchange_TokenExchangeFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAdapter := NewMockAdapter(ctrl)
	svc := NewService(mockAdapter, "client-id", "client-secret")

	dto := &models.GoogleLoginRequestDTO{Code: "bad-code", RedirectURI: "http://localhost"}

	mockAdapter.EXPECT().ExchangeToken(gomock.Any(), "bad-code", "client-id", "client-secret", "http://localhost").
		Return(nil, errors.New("invalid code"))

	_, err := svc.Exchange(dto)
	if err == nil {
		t.Fatal("expected error for failed token exchange")
	}
}

func TestExchange_GetUserInfoFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAdapter := NewMockAdapter(ctrl)
	svc := NewService(mockAdapter, "client-id", "client-secret")

	dto := &models.GoogleLoginRequestDTO{Code: "code", RedirectURI: "http://localhost"}

	mockAdapter.EXPECT().ExchangeToken(gomock.Any(), "code", "client-id", "client-secret", "http://localhost").
		Return(&TokenResponse{AccessToken: "token"}, nil)

	mockAdapter.EXPECT().GetUserInfo(gomock.Any(), "token").
		Return(nil, errors.New("userinfo failed"))

	_, err := svc.Exchange(dto)
	if err == nil {
		t.Fatal("expected error for failed user info")
	}
}
