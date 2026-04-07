//go:generate mockgen -source=google_oauth_service.go -destination=google_oauth_service_mock.go -package=google_oauth

package google_oauth

import (
	"context"
	"fmt"
)

type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
	Picture  string
}

type Service interface {
	Exchange(ctx context.Context, code, redirectURI string) (*GoogleUserInfo, error)
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

func (s *service) Exchange(ctx context.Context, code, redirectURI string) (*GoogleUserInfo, error) {
	tokenResp, err := s.adapter.ExchangeToken(ctx, code, s.clientID, s.clientSecret, redirectURI)
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
