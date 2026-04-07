package controllers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
	svcerr "github.com/lauralesteves/copa-guru-backend/internal/shared/services"
	"go.uber.org/mock/gomock"
)

func setupController(t *testing.T) (*gomock.Controller, *services.MockAuthService, *AuthController) {
	ctrl := gomock.NewController(t)
	mockAuth := services.NewMockAuthService(ctrl)
	controller := NewAuthController(mockAuth)
	return ctrl, mockAuth, controller
}

func authedRequest(userID, body string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		Body: body,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{"userId": userID},
		},
	}
}

// --- GoogleLogin ---

func TestGoogleLogin_Success(t *testing.T) {
	ctrl, mockAuth, controller := setupController(t)
	defer ctrl.Finish()

	now := time.Now()
	mockAuth.EXPECT().LoginWithGoogle(gomock.Any()).Return(&models.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		User:         &models.UserDTO{ID: "user-1", Email: "a@b.com", CreatedAt: now},
	}, nil)

	req := events.APIGatewayProxyRequest{
		Body: `{"code":"auth-code","redirectUri":"http://localhost:5173"}`,
	}
	resp, err := controller.GoogleLogin(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	var body models.LoginResponse
	json.Unmarshal([]byte(resp.Body), &body)
	if body.AccessToken != "access-token" {
		t.Errorf("AccessToken = %q, want %q", body.AccessToken, "access-token")
	}
}

func TestGoogleLogin_InvalidBody(t *testing.T) {
	ctrl, _, controller := setupController(t)
	defer ctrl.Finish()

	req := events.APIGatewayProxyRequest{Body: "not json"}
	resp, _ := controller.GoogleLogin(req)
	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestGoogleLogin_MissingFields(t *testing.T) {
	ctrl, _, controller := setupController(t)
	defer ctrl.Finish()

	req := events.APIGatewayProxyRequest{Body: `{"code":""}`}
	resp, _ := controller.GoogleLogin(req)
	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestGoogleLogin_ServiceError(t *testing.T) {
	ctrl, mockAuth, controller := setupController(t)
	defer ctrl.Finish()

	mockAuth.EXPECT().LoginWithGoogle(gomock.Any()).Return(nil, svcerr.NewInternalError("failed", nil))

	req := events.APIGatewayProxyRequest{Body: `{"code":"bad","redirectUri":"http://x"}`}
	resp, _ := controller.GoogleLogin(req)
	if resp.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", resp.StatusCode)
	}
}

// --- Refresh ---

func TestRefresh_Success(t *testing.T) {
	ctrl, mockAuth, controller := setupController(t)
	defer ctrl.Finish()

	mockAuth.EXPECT().RefreshTokens("old-token").Return(&models.TokenResponse{
		AccessToken:  "new-access",
		RefreshToken: "new-refresh",
	}, nil)

	req := events.APIGatewayProxyRequest{Body: `{"refreshToken":"old-token"}`}
	resp, _ := controller.Refresh(req)
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	var body models.TokenResponse
	json.Unmarshal([]byte(resp.Body), &body)
	if body.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want %q", body.AccessToken, "new-access")
	}
}

func TestRefresh_InvalidBody(t *testing.T) {
	ctrl, _, controller := setupController(t)
	defer ctrl.Finish()

	req := events.APIGatewayProxyRequest{Body: "bad"}
	resp, _ := controller.Refresh(req)
	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestRefresh_MissingToken(t *testing.T) {
	ctrl, _, controller := setupController(t)
	defer ctrl.Finish()

	req := events.APIGatewayProxyRequest{Body: `{"refreshToken":""}`}
	resp, _ := controller.Refresh(req)
	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestRefresh_Unauthorized(t *testing.T) {
	ctrl, mockAuth, controller := setupController(t)
	defer ctrl.Finish()

	mockAuth.EXPECT().RefreshTokens("expired").Return(nil, svcerr.NewUnauthorizedError("expired"))

	req := events.APIGatewayProxyRequest{Body: `{"refreshToken":"expired"}`}
	resp, _ := controller.Refresh(req)
	if resp.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", resp.StatusCode)
	}
}

// --- Logout ---

func TestLogout_Success(t *testing.T) {
	ctrl, mockAuth, controller := setupController(t)
	defer ctrl.Finish()

	mockAuth.EXPECT().Logout("user-42").Return(nil)

	resp, _ := controller.Logout(authedRequest("user-42", ""))
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
}

func TestLogout_NoAuth(t *testing.T) {
	ctrl, _, controller := setupController(t)
	defer ctrl.Finish()

	req := events.APIGatewayProxyRequest{}
	resp, _ := controller.Logout(req)
	if resp.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", resp.StatusCode)
	}
}

// --- Me ---

func TestMe_Success(t *testing.T) {
	ctrl, mockAuth, controller := setupController(t)
	defer ctrl.Finish()

	mockAuth.EXPECT().GetMe("user-42").Return(&models.MeDTO{
		User: &models.UserDTO{ID: "user-42", Email: "me@copa.guru"},
	}, nil)

	resp, _ := controller.Me(authedRequest("user-42", ""))
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	var body models.MeDTO
	json.Unmarshal([]byte(resp.Body), &body)
	if body.User.Email != "me@copa.guru" {
		t.Errorf("Email = %q, want %q", body.User.Email, "me@copa.guru")
	}
}

func TestMe_NoAuth(t *testing.T) {
	ctrl, _, controller := setupController(t)
	defer ctrl.Finish()

	req := events.APIGatewayProxyRequest{}
	resp, _ := controller.Me(req)
	if resp.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", resp.StatusCode)
	}
}

func TestMe_NotFound(t *testing.T) {
	ctrl, mockAuth, controller := setupController(t)
	defer ctrl.Finish()

	mockAuth.EXPECT().GetMe("user-99").Return(nil, svcerr.NewNotFoundError("user not found"))

	resp, _ := controller.Me(authedRequest("user-99", ""))
	if resp.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", resp.StatusCode)
	}
}
