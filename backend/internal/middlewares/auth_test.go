package middlewares

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

func validToken(secret string) string {
	claims := jwt.MapClaims{
		"userId": "user-123",
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func expiredToken(secret string) string {
	claims := jwt.MapClaims{
		"userId": "user-123",
		"exp":    time.Now().Add(-1 * time.Hour).Unix(),
		"iat":    time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func authRequest(headers map[string]string) events.APIGatewayCustomAuthorizerRequestTypeRequest {
	return events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers:   headers,
		MethodArn: "arn:aws:execute-api:us-east-1:123456789:api/GET/resource",
	}
}

func TestAuthorizeBearer_Valid(t *testing.T) {
	secret := "copa-guru-local-jwt-secret-key"
	t.Setenv("AUTH_JWT_SECRET", secret)

	req := authRequest(map[string]string{"Authorization": "Bearer " + validToken(secret)})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Allow" {
		t.Errorf("Effect = %q, want Allow", resp.PolicyDocument.Statement[0].Effect)
	}
}

func TestAuthorizeBearer_NoHeader(t *testing.T) {
	req := authRequest(map[string]string{})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Errorf("Effect = %q, want Deny", resp.PolicyDocument.Statement[0].Effect)
	}
}

func TestAuthorizeBearer_NoBearerPrefix(t *testing.T) {
	req := authRequest(map[string]string{"Authorization": "Basic abc123"})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Errorf("Effect = %q, want Deny", resp.PolicyDocument.Statement[0].Effect)
	}
}

func TestAuthorizeBearer_ExpiredToken(t *testing.T) {
	secret := "copa-guru-local-jwt-secret-key"
	t.Setenv("AUTH_JWT_SECRET", secret)

	req := authRequest(map[string]string{"Authorization": "Bearer " + expiredToken(secret)})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Errorf("Effect = %q, want Deny", resp.PolicyDocument.Statement[0].Effect)
	}
}

func TestAuthorizeBearer_InvalidToken(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "copa-guru-local-jwt-secret-key")

	req := authRequest(map[string]string{"Authorization": "Bearer not-a-real-token"})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Errorf("Effect = %q, want Deny", resp.PolicyDocument.Statement[0].Effect)
	}
}

func TestAuthorizeBearer_WrongSecret(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "correct-secret")

	req := authRequest(map[string]string{"Authorization": "Bearer " + validToken("wrong-secret")})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Errorf("Effect = %q, want Deny", resp.PolicyDocument.Statement[0].Effect)
	}
}

func TestAuthorizeBearer_EmptySecret(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "")

	req := authRequest(map[string]string{"Authorization": "Bearer " + validToken("any")})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Errorf("Effect = %q, want Deny", resp.PolicyDocument.Statement[0].Effect)
	}
}

func TestAuthorizeBearer_XAuthHeader(t *testing.T) {
	secret := "copa-guru-local-jwt-secret-key"
	t.Setenv("AUTH_JWT_SECRET", secret)

	req := authRequest(map[string]string{"X-Auth": "Bearer " + validToken(secret)})
	resp, err := AuthorizeBearer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Allow" {
		t.Errorf("Effect = %q, want Allow", resp.PolicyDocument.Statement[0].Effect)
	}
}

// --- Helper unit tests ---

func TestGetAuthHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{"Authorization", map[string]string{"Authorization": "Bearer abc"}, "Bearer abc"},
		{"authorization lowercase", map[string]string{"authorization": "Bearer abc"}, "Bearer abc"},
		{"X-Auth", map[string]string{"X-Auth": "Bearer abc"}, "Bearer abc"},
		{"x-auth lowercase", map[string]string{"x-auth": "Bearer abc"}, "Bearer abc"},
		{"empty", map[string]string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := events.APIGatewayCustomAuthorizerRequestTypeRequest{Headers: tt.headers}
			got := getAuthHeader(req)
			if got != tt.want {
				t.Errorf("getAuthHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		auth      string
		wantOk    bool
		wantToken string
	}{
		{"valid bearer", "Bearer my-token", true, "my-token"},
		{"no prefix", "my-token", false, ""},
		{"basic auth", "Basic abc", false, ""},
		{"empty", "", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, token := getBearerToken(tt.auth)
			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if token != tt.wantToken {
				t.Errorf("token = %q, want %q", token, tt.wantToken)
			}
		})
	}
}
