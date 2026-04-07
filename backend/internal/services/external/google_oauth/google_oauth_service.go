//go:generate mockgen -source=google_oauth_service.go -destination=google_oauth_service_mock.go -package=google_oauth

package google_oauth

import (
	"context"
	"fmt"

	"github.com/lauralesteves/copa-guru-backend/internal/models"
)

type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
	Picture  string
}

type Service interface {
	Exchange(dto *models.LoginRequestDTO) (*GoogleUserInfo, error)
}

type service struct {
	adapter      Adapter
	clientID     string
	clientSecret string
}

func NewService(adapter Adapter, clientID, clientSecret string) Service {
	return &service{
		adapter:      adapter,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *service) Exchange(dto *models.LoginRequestDTO) (*GoogleUserInfo, error) {
	ctx := context.Background()

	tokenResp, err := s.adapter.ExchangeToken(ctx, dto.Code, s.clientID, s.clientSecret, dto.RedirectURI)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	userInfo, err := s.adapter.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &GoogleUserInfo{
		GoogleID: userInfo.Sub,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
	}, nil
}
