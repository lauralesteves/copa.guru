package controllers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
)

func parseGoogleLoginRequest(request events.APIGatewayProxyRequest) (*models.GoogleLoginRequestDTO, events.APIGatewayProxyResponse) {
	return parseAndValidate[models.GoogleLoginRequestDTO](request)
}

func parseRefreshRequest(request events.APIGatewayProxyRequest) (*models.RefreshRequestDTO, events.APIGatewayProxyResponse) {
	return parseAndValidate[models.RefreshRequestDTO](request)
}
