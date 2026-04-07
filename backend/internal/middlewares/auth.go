package middlewares

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// AuthorizeBearer accepts only bearer Token
func AuthorizeBearer(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	hasBearerToken, bearerToken := getBearerToken(getAuthHeader(request))

	if !hasBearerToken {
		return unauthorized(), nil
	}

	if !validateBearerToken(bearerToken) {
		return unauthorized(), nil
	}

	return allowPolicy(request.MethodArn), nil
}
