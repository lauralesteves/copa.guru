package controllers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
	"github.com/lauralesteves/copa-guru-backend/internal/shared/helpers"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) GoogleLogin(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dto, validationResponse := parseGoogleLoginRequest(request)
	if dto == nil {
		return validationResponse, nil
	}

	resp, err := c.authService.LoginWithGoogle(dto)
	if err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(resp, http.StatusOK), nil
}

func (c *AuthController) Refresh(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dto, validationResponse := parseRefreshRequest(request)
	if dto == nil {
		return validationResponse, nil
	}

	resp, err := c.authService.RefreshTokens(dto.RefreshToken)
	if err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(resp, http.StatusOK), nil
}

func (c *AuthController) Logout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := helpers.GetId(request)
	if id == "" {
		return helpers.UnauthorizedResponse(), nil
	}

	if err := c.authService.Logout(id); err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(map[string]string{"message": "logged out"}, http.StatusOK), nil
}

func (c *AuthController) Me(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := helpers.GetId(request)
	if id == "" {
		return helpers.UnauthorizedResponse(), nil
	}

	resp, err := c.authService.GetMe(id)
	if err != nil {
		return helpers.HandleServiceError(err), nil
	}

	return helpers.SuccessResponse(resp, http.StatusOK), nil
}
