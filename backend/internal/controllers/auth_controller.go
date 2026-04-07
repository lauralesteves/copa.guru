package controllers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
	"github.com/lauralesteves/copa-guru-backend/internal/shared/helpers"
	svcerr "github.com/lauralesteves/copa-guru-backend/internal/shared/services"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) GoogleLogin(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.GoogleLoginRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		slog.Warn("invalid google login request body", "error", err)
		return helpers.BadRequestResponse("invalid request body"), nil
	}

	if err := req.Validate(); err != nil {
		return helpers.ServiceErrorResponse(svcerr.NewValidationError(err.Error())), nil
	}

	resp, err := c.authService.LoginWithGoogle(context.Background(), req.Code, req.RedirectURI)
	if err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(resp, http.StatusOK), nil
}

func (c *AuthController) Refresh(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.RefreshRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		slog.Warn("invalid refresh request body", "error", err)
		return helpers.BadRequestResponse("invalid request body"), nil
	}

	if err := req.Validate(); err != nil {
		return helpers.ServiceErrorResponse(svcerr.NewValidationError(err.Error())), nil
	}

	resp, err := c.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(resp, http.StatusOK), nil
}

func (c *AuthController) Logout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := getUserIDFromContext(request)
	if userID == "" {
		return helpers.UnauthorizedResponse(), nil
	}

	if err := c.authService.Logout(userID); err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(map[string]string{"message": "logged out"}, http.StatusOK), nil
}

func (c *AuthController) Me(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := getUserIDFromContext(request)
	if userID == "" {
		return helpers.UnauthorizedResponse(), nil
	}

	resp, err := c.authService.GetMe(userID)
	if err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(resp, http.StatusOK), nil
}

func getUserIDFromContext(request events.APIGatewayProxyRequest) string {
	if auth, ok := request.RequestContext.Authorizer["userId"]; ok {
		if userID, ok := auth.(string); ok {
			return userID
		}
	}
	return ""
}
