package controllers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
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

	user, err := c.authService.LoginWithGoogle(dto)
	if err != nil {
		return middlewares.HandleServiceError(err), nil
	}

	return middlewares.SuccessResponse(user.ToLoginResponseDTO(), http.StatusOK), nil
}

func (c *AuthController) Refresh(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dto, validationResponse := parseRefreshRequest(request)
	if dto == nil {
		return validationResponse, nil
	}

	auth, err := c.authService.RefreshTokens(dto.RefreshToken)
	if err != nil {
		return middlewares.HandleServiceError(err), nil
	}

	return middlewares.SuccessResponse(auth.ToRefreshResponseDTO(), http.StatusOK), nil
}

func (c *AuthController) Logout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := middlewares.GetId(request)
	if id == "" {
		return middlewares.UnauthorizedResponse(), nil
	}

	if err := c.authService.Logout(id); err != nil {
		return middlewares.HandleServiceError(err), nil
	}

	return middlewares.SuccessResponse(map[string]string{"message": "logged out"}, http.StatusOK), nil
}

func (c *AuthController) Me(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := middlewares.GetId(request)
	if id == "" {
		return middlewares.UnauthorizedResponse(), nil
	}

	user, err := c.authService.GetMe(id)
	if err != nil {
		return middlewares.HandleServiceError(err), nil
	}

	return middlewares.SuccessResponse(user.ToDTO(), http.StatusOK), nil
}
