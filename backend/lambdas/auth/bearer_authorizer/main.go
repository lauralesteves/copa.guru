package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lauralesteves/copa-guru-backend/internal/shared/config"
)

func handler(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.Headers["authorization"]
	if token == "" {
		token = request.Headers["Authorization"]
	}

	if !strings.HasPrefix(token, "Bearer ") {
		slog.Warn("missing or invalid authorization header")
		return generatePolicy("", "Deny", request.MethodArn), nil
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")

	userID, err := validateToken(tokenString, config.GetJWTSecret())
	if err != nil {
		slog.Warn("token validation failed", "error", err)
		return generatePolicy("", "Deny", request.MethodArn), nil
	}

	return generatePolicy(userID, "Allow", request.MethodArn), nil
}

func validateToken(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["userId"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("missing userId in claims")
	}

	return userID, nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	if principalID == "" {
		principalID = "unauthorized"
	}

	resp := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalID,
	}

	if effect != "" && resource != "" {
		resp.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	if effect == "Allow" {
		resp.Context = map[string]interface{}{
			"userId": principalID,
		}
	}

	return resp
}

func main() {
	lambda.Start(handler)
}
