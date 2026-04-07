package middlewares

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lauralesteves/copa-guru-backend/internal/config"
)

const bearerPrefix = "Bearer "

func allowPolicy(methodArn string) events.APIGatewayCustomAuthorizerResponse {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{methodArn},
				},
			},
		},
	}
}

func unauthorized() events.APIGatewayCustomAuthorizerResponse {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Deny",
					Resource: []string{"*"},
				},
			},
		},
	}
}

func getAuthHeader(request events.APIGatewayCustomAuthorizerRequestTypeRequest) string {
	for _, key := range []string{"X-Auth", "x-auth", "Authorization", "authorization"} {
		if auth := request.Headers[key]; auth != "" {
			return auth
		}
	}
	return ""
}

func getBearerToken(auth string) (bool, string) {
	if !strings.HasPrefix(auth, bearerPrefix) {
		return false, ""
	}

	return true, auth[len(bearerPrefix):]
}

func validateBearerToken(tokenString string) bool {
	secret := config.GetJWTSecret()
	if secret == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	return err == nil && token.Valid
}
