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

// GoogleLogin handles POST /auth/google requests.
// Authenticates a user via Google OAuth authorization code exchange.
//
// Body Parameters:
//   - code: Google OAuth authorization code (required)
//   - redirectUri: OAuth redirect URI used in the frontend flow (required)
//
// Responses:
//   - 200: LoginResponseDTO with access token, refresh token, and user data
//   - 400: Invalid request body or missing fields
//   - 500: Google exchange or database failure
//
// Authentication: None
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

// Refresh handles POST /auth/refresh requests.
// Rotates access and refresh tokens for an authenticated session.
//
// Body Parameters:
//   - refreshToken: current refresh token (required)
//
// Responses:
//   - 200: RefreshResponseDTO with new access and refresh tokens
//   - 400: Invalid request body or missing fields
//   - 401: Invalid or expired refresh token
//
// Authentication: None
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

// Logout handles POST /auth/logout requests.
// Invalidates the user's refresh token, ending the session.
//
// Responses:
//   - 200: { "message": "logged out" }
//   - 401: Missing or invalid authorization
//   - 500: Database failure
//
// Authentication: Bearer token
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

// Me handles GET /auth/me requests.
// Returns the authenticated user's profile.
//
// Responses:
//   - 200: UserDTO with id, email, name, picture, lastLoginAt, createdAt
//   - 401: Missing or invalid authorization
//   - 404: User not found
//
// Authentication: Bearer token
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
