//go:generate mockgen -source=google_oauth.go -destination=google_oauth_mock.go -package=services

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
	Picture  string
}

type GoogleOAuthClient interface {
	Exchange(ctx context.Context, code, redirectURI string) (*GoogleUserInfo, error)
}

type googleOAuthClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewGoogleOAuthClient(clientID, clientSecret string) GoogleOAuthClient {
	return &googleOAuthClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   http.DefaultClient,
	}
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type googleUserInfoResponse struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (c *googleOAuthClient) Exchange(ctx context.Context, code, redirectURI string) (*GoogleUserInfo, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp googleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	userInfo, err := c.getUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (c *googleOAuthClient) getUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo failed (%d)", resp.StatusCode)
	}

	var info googleUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo: %w", err)
	}

	return &GoogleUserInfo{
		GoogleID: info.Sub,
		Email:    info.Email,
		Name:     info.Name,
		Picture:  info.Picture,
	}, nil
}
