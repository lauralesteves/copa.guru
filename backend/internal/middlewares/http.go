package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
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

	id = request.PathParameters["userId"]
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

func GetQueryString(request events.APIGatewayProxyRequest, param string) string {
	return request.QueryStringParameters[param]
}

func GetQueryStringInt(request events.APIGatewayProxyRequest, param string, defaultValue int64) int64 {
	result := defaultValue

	if value := GetQueryString(request, param); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil && parsed > 0 && parsed <= 100 {
			result = parsed
		}
	}

	return result
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

func ServiceErrorResponse(svcErr *ServiceError) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": svcErr.Message})
	return events.APIGatewayProxyResponse{
		StatusCode: int(svcErr.Code),
		Body:       string(body),
		Headers:    CorsHeaders,
	}
}

func HandleServiceError(err error) events.APIGatewayProxyResponse {
	var svcErr *ServiceError
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
