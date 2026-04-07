package config

import (
	"os"
	"testing"
)

func TestGetJWTSecret_Default(t *testing.T) {
	os.Unsetenv("AUTH_JWT_SECRET")
	if got := GetJWTSecret(); got != "copa-guru-local-jwt-secret-key" {
		t.Errorf("GetJWTSecret() = %q, want %q", got, "copa-guru-local-jwt-secret-key")
	}
}

func TestGetJWTSecret_FromEnv(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "prod-secret")
	if got := GetJWTSecret(); got != "prod-secret" {
		t.Errorf("GetJWTSecret() = %q, want %q", got, "prod-secret")
	}
}

func TestGetGoogleClientID_Empty(t *testing.T) {
	os.Unsetenv("GOOGLE_CLIENT_ID")
	if got := GetGoogleClientID(); got != "" {
		t.Errorf("GetGoogleClientID() = %q, want empty", got)
	}
}

func TestGetGoogleClientID_FromEnv(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id-123")
	if got := GetGoogleClientID(); got != "client-id-123" {
		t.Errorf("GetGoogleClientID() = %q, want %q", got, "client-id-123")
	}
}

func TestGetGoogleClientSecret_Empty(t *testing.T) {
	os.Unsetenv("GOOGLE_CLIENT_SECRET")
	if got := GetGoogleClientSecret(); got != "" {
		t.Errorf("GetGoogleClientSecret() = %q, want empty", got)
	}
}

func TestGetGoogleClientSecret_FromEnv(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_SECRET", "secret-456")
	if got := GetGoogleClientSecret(); got != "secret-456" {
		t.Errorf("GetGoogleClientSecret() = %q, want %q", got, "secret-456")
	}
}
