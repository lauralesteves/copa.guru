//go:generate mockgen -source=auth_service.go -destination=auth_service_mock.go -package=services

package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	svcerr "github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"github.com/lauralesteves/copa-guru-backend/internal/services/external/google_oauth"
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
	LoginWithGoogle(dto *models.GoogleLoginRequestDTO) (*models.User, error)
	RefreshTokens(refreshToken string) (*models.Auth, error)
	Logout(userID string) error
	GetMe(userID string) (*models.User, error)
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

func (s *authService) LoginWithGoogle(dto *models.GoogleLoginRequestDTO) (*models.User, error) {
	googleUser, err := s.googleOAuth.Exchange(dto)
	if err != nil {
		slog.Error("google oauth exchange failed", "error", err)
		return nil, svcerr.NewInternalError("failed to authenticate with Google", err)
	}

	now := time.Now()

	auth := models.Auth{
		GoogleID:    googleUser.GoogleID,
		LastLoginAt: &now,
	}

	user := &models.User{
		Strategy: models.StrategyGoogle,

		Name:    googleUser.Name,
		Email:   googleUser.Email,
		Picture: googleUser.Picture,

		Auth: auth,
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
	user.Auth.AccessToken = accessToken

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		slog.Error("failed to generate refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to generate refresh token", err)
	}
	user.Auth.RefreshToken = refreshToken

	if err = s.userRepo.UpdateRefreshToken(user.ID, refreshToken, now.Add(RefreshTokenTTL)); err != nil {
		slog.Error("failed to save refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to save refresh token", err)
	}

	return user, nil
}

func (s *authService) RefreshTokens(refreshToken string) (*models.Auth, error) {
	user, err := s.userRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		slog.Error("failed to find user by refresh token", "error", err)
		return nil, svcerr.NewInternalError("failed to refresh tokens", err)
	}
	if user == nil {
		return nil, svcerr.NewUnauthorizedError("invalid refresh token")
	}

	auth := user.Auth

	if auth.RefreshTokenExpiresAt != nil && auth.RefreshTokenExpiresAt.Before(time.Now()) {
		_ = s.userRepo.InvalidateRefreshToken(user.ID)
		return nil, svcerr.NewUnauthorizedError("refresh token expired")
	}

	accessToken, err := s.generateAccessToken(user.ID.Hex())
	if err != nil {
		slog.Error("failed to generate access token on refresh", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to generate access token", err)
	}
	auth.AccessToken = accessToken

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		slog.Error("failed to generate new refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to generate refresh token", err)
	}
	auth.RefreshToken = newRefreshToken

	if err = s.userRepo.UpdateRefreshToken(user.ID, newRefreshToken, time.Now().Add(RefreshTokenTTL)); err != nil {
		slog.Error("failed to rotate refresh token", "userId", user.ID.Hex(), "error", err)
		return nil, svcerr.NewInternalError("failed to save refresh token", err)
	}

	return &auth, nil
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

func (s *authService) GetMe(userID string) (*models.User, error) {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, svcerr.NewValidationError("invalid user ID")
	}

	user, err := s.userRepo.Get(id)
	if err != nil {
		slog.Error("failed to find user", "userId", userID, "error", err)
		return nil, svcerr.NewInternalError("failed to get user", err)
	}
	if user == nil {
		return nil, svcerr.NewNotFoundError("user not found")
	}

	return user, nil
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
