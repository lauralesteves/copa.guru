package controllers

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/lauralesteves/copa-guru-backend/internal/middlewares"
)

var validate = validator.New()

func parseAndValidate[T any](request events.APIGatewayProxyRequest) (*T, events.APIGatewayProxyResponse) {
	var dto T
	if err := json.Unmarshal([]byte(request.Body), &dto); err != nil {
		slog.Warn("invalid request body", "error", err)
		return nil, middlewares.BadRequestResponse([]string{err.Error()})
	}

	if err := validate.Struct(dto); err != nil {
		var errs []string
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, formatValidationError(e))
		}
		slog.Warn("validation failed", "errors", errs)
		return nil, middlewares.BadRequestResponse(errs)
	}

	return &dto, events.APIGatewayProxyResponse{}
}

func formatValidationError(e validator.FieldError) string {
	field := e.Namespace()
	switch e.Tag() {
	case "required", "required_if":
		return "missing \"" + field + "\" in body"
	case "oneof":
		return "invalid \"" + field + "\" in body, must be one of: " + e.Param()
	default:
		return "invalid \"" + field + "\" in body"
	}
}
