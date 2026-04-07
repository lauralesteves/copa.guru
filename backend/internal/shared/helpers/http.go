package helpers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lauralesteves/copa-guru-backend/internal/shared/services"
)

var CorsHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Headers": "Content-Type,Authorization",
	"Access-Control-Allow-Methods": "GET,POST,PUT,PATCH,DELETE,OPTIONS",
}

func GetId(request events.APIGatewayProxyRequest) string {
	var id string

	// Get from path
	id = request.PathParameters["id"]
	if id != "" {
		return id
	}

	// Get from authorizer
	if auth, ok := request.RequestContext.Authorizer["userId"]; ok {
		if id, ok = auth.(string); ok {
			return id
		}
	}
	return id
}

func SuccessResponse(data interface{}, statusCode int) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    CorsHeaders,
	}
	if data != nil {
		body, _ := json.Marshal(data)
		resp.Body = string(body)
	}
	return resp
}

func BadRequestResponse(errs []string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]interface{}{"errors": errs})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Headers:    CorsHeaders,
		Body:       string(body),
	}
}

func NotFoundResponse(entity string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": entity + " not found"})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       string(body),
		Headers:    CorsHeaders,
	}
}

func UnauthorizedResponse() events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": "unauthorized"})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       string(body),
		Headers:    CorsHeaders,
	}
}

func ForbiddenResponse(message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusForbidden,
		Body:       string(body),
		Headers:    CorsHeaders,
	}
}

func ServiceErrorResponse(svcErr *services.ServiceError) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": svcErr.Message})
	return events.APIGatewayProxyResponse{
		StatusCode: int(svcErr.Code),
		Body:       string(body),
		Headers:    CorsHeaders,
	}
}

func HandleServiceError(err error) events.APIGatewayProxyResponse {
	var svcErr *services.ServiceError
	if errors.As(err, &svcErr) {
		return ServiceErrorResponse(svcErr)
	}
	body, _ := json.Marshal(map[string]string{"error": "internal server error"})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(body),
		Headers:    CorsHeaders,
	}
}
