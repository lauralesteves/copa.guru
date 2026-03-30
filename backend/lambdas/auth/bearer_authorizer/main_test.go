package main

import (
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key"

func createToken(claims jwt.MapClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(secret))
	return s
}

func TestValidateToken_Valid(t *testing.T) {
	token := createToken(jwt.MapClaims{"userId": "abc123"}, testSecret)
	userID, err := validateToken(token, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "abc123" {
		t.Errorf("userID = %q, want %q", userID, "abc123")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	token := createToken(jwt.MapClaims{
		"userId": "abc123",
		"exp":    time.Now().Add(-1 * time.Hour).Unix(),
	}, testSecret)
	_, err := validateToken(token, testSecret)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token := createToken(jwt.MapClaims{"userId": "abc123"}, "wrong-secret")
	_, err := validateToken(token, testSecret)
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestValidateToken_Malformed(t *testing.T) {
	_, err := validateToken("not-a-jwt-token", testSecret)
	if err == nil {
		t.Error("expected error for malformed token")
	}
}

func TestValidateToken_MissingUserId(t *testing.T) {
	token := createToken(jwt.MapClaims{"sub": "abc123"}, testSecret)
	_, err := validateToken(token, testSecret)
	if err == nil {
		t.Error("expected error for missing userId")
	}
}

func TestValidateToken_WrongSigningMethod(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"userId": "abc123"})
	s, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	_, err := validateToken(s, testSecret)
	if err == nil {
		t.Error("expected error for none signing method")
	}
}

func TestHandler_ValidToken(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", testSecret)
	token := createToken(jwt.MapClaims{"userId": "user-42"}, testSecret)
	req := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers:   map[string]string{"Authorization": "Bearer " + token},
		MethodArn: "arn:aws:execute-api:us-east-1:123:api/GET/test",
	}

	resp, err := handler(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PrincipalID != "user-42" {
		t.Errorf("PrincipalID = %q, want %q", resp.PrincipalID, "user-42")
	}
	if resp.PolicyDocument.Statement[0].Effect != "Allow" {
		t.Error("expected Allow effect")
	}
	if resp.Context["userId"] != "user-42" {
		t.Errorf("context userId = %v, want %q", resp.Context["userId"], "user-42")
	}
}

func TestHandler_MissingToken(t *testing.T) {
	req := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers:   map[string]string{},
		MethodArn: "arn:aws:execute-api:us-east-1:123:api/GET/test",
	}

	resp, err := handler(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Error("expected Deny effect")
	}
}

func TestHandler_InvalidToken(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", testSecret)
	req := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers:   map[string]string{"Authorization": "Bearer invalid-token"},
		MethodArn: "arn:aws:execute-api:us-east-1:123:api/GET/test",
	}

	resp, err := handler(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Error("expected Deny effect")
	}
}

func TestHandler_NoBearerPrefix(t *testing.T) {
	req := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers:   map[string]string{"Authorization": "Basic abc123"},
		MethodArn: "arn:aws:execute-api:us-east-1:123:api/GET/test",
	}

	resp, err := handler(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PolicyDocument.Statement[0].Effect != "Deny" {
		t.Error("expected Deny effect")
	}
}

func TestGeneratePolicy_Allow(t *testing.T) {
	resp := generatePolicy("user-1", "Allow", "arn:test")
	if resp.PrincipalID != "user-1" {
		t.Errorf("PrincipalID = %q, want %q", resp.PrincipalID, "user-1")
	}
	if resp.Context["userId"] != "user-1" {
		t.Error("context should contain userId on Allow")
	}
}

func TestGeneratePolicy_Deny(t *testing.T) {
	resp := generatePolicy("", "Deny", "arn:test")
	if resp.PrincipalID != "unauthorized" {
		t.Errorf("PrincipalID = %q, want %q", resp.PrincipalID, "unauthorized")
	}
	if resp.Context != nil {
		t.Error("context should be nil on Deny")
	}
}
