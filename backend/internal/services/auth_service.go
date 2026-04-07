//go:generate mockgen -source=auth_service.go -destination=auth_service_mock.go -package=services

package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"github.com/lauralesteves/copa-guru-backend/internal/services/external/google_oauth"
	svcerr "github.com/lauralesteves/copa-guru-backend/internal/shared/services"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	AccessTokenTTL  = 1 * time.Hour
	RefreshTokenTTL = 30 * 24 * time.Hour
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type AuthService interface {
	LoginWithGoogle(ctx context.Context, code, redirectURI string) (*models.LoginResponse, error)
	RefreshTokens(refreshToken string) (*models.TokenResponse, error)
	Logout(userID string) error
	GetMe(userID string) (*models.MeDTO, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	jwtSecret   string
	googleOAuth google_oauth.Service
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string, googleOAuth google_oauth.Service) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		googleOAuth: googleOAuth,
	}
}

func (s *authService) generateAccessToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *authService) validateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, nil
}

func (s *authService) LoginWithGoogle(ctx context.Context, code, redirectURI string) (*models.LoginResponse, error) {
	googleUser, err := s.googleOAuth.Exchange(ctx, code, redirectURI)
	if err != nil {
		slog.Error("google oauth exchange failed", "error", err)
		return nil, svcerr.NewInternalError("failed to authenticate with Google", err)
	}

	now := time.Now()
	user := &models.User{
		Strategy:    models.StrategyGoogle,
		GoogleID:    googleUser.GoogleID,
		Email:       googleUser.Email,
		Name:        googleUser.Name,
		Picture:     googleUser.Picture,
		LastLoginAt: &now,
	}

	user, err = s.userRepo.Upsert(user)
	if err != nil {
		slog.Error("failed to upsert user", "email", googleUser.Email, "error", err)
		return nil, svcerr.NewInternalError("failed to save user", err)
	}

	accessToken, err := s.generateAccessToken(user.ID.Hex())
	if err != nil {
		slog.Error("failed to generate access token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to generate access token", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		slog.Error("failed to generate refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to generate refresh token", err)
	}

	if err = s.userRepo.UpdateRefreshToken(user.ID, refreshToken, now.Add(RefreshTokenTTL)); err != nil {
		slog.Error("failed to save refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to save refresh token", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToDTO(),
	}, nil
}

func (s *authService) RefreshTokens(refreshToken string) (*models.TokenResponse, error) {
	user, err := s.userRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		slog.Error("failed to find user by refresh token", "error", err)
		return nil, svcerr.NewInternalError("failed to refresh tokens", err)
	}
	if user == nil {
		return nil, svcerr.NewUnauthorizedError("invalid refresh token")
	}

	if user.RefreshTokenExpiresAt != nil && user.RefreshTokenExpiresAt.Before(time.Now()) {
		_ = s.userRepo.InvalidateRefreshToken(user.ID)
		return nil, svcerr.NewUnauthorizedError("refresh token expired")
	}

	accessToken, err := s.generateAccessToken(user.ID.Hex())
	if err != nil {
		slog.Error("failed to generate access token on refresh", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to generate access token", err)
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		slog.Error("failed to generate new refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to generate refresh token", err)
	}

	if err := s.userRepo.UpdateRefreshToken(user.ID, newRefreshToken, time.Now().Add(RefreshTokenTTL)); err != nil {
		slog.Error("failed to rotate refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to save refresh token", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) Logout(userID string) error {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return svcerr.NewValidationError("invalid user ID")
	}

	if err := s.userRepo.InvalidateRefreshToken(id); err != nil {
		slog.Error("failed to invalidate refresh token", "userId", userID, "error", err)
		return svcerr.NewInternalError("failed to logout", err)
	}

	return nil
}

func (s *authService) GetMe(userID string) (*models.MeDTO, error) {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, svcerr.NewValidationError("invalid user ID")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		slog.Error("failed to find user", "userId", userID, "error", err)
		return nil, svcerr.NewInternalError("failed to get user", err)
	}
	if user == nil {
		return nil, svcerr.NewNotFoundError("user not found")
	}

	return &models.MeDTO{
		User: user.ToDTO(),
	}, nil
}
