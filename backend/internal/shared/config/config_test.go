package config

import (
	"os"
	"testing"
)

func TestGetDatabaseName_Default(t *testing.T) {
	os.Unsetenv("MONGODB_DATABASE_NAME")
	if got := GetDatabaseName(); got != "copa-guru" {
		t.Errorf("GetDatabaseName() = %q, want %q", got, "copa-guru")
	}
}

func TestGetDatabaseName_FromEnv(t *testing.T) {
	t.Setenv("MONGODB_DATABASE_NAME", "custom-db")
	if got := GetDatabaseName(); got != "custom-db" {
		t.Errorf("GetDatabaseName() = %q, want %q", got, "custom-db")
	}
}

func TestGetDatabaseDSN_Default(t *testing.T) {
	os.Unsetenv("MONGODB_DSN")
	if got := GetDatabaseDSN(); got != "mongodb://localhost:27017" {
		t.Errorf("GetDatabaseDSN() = %q, want %q", got, "mongodb://localhost:27017")
	}
}

func TestGetDatabaseDSN_FromEnv(t *testing.T) {
	t.Setenv("MONGODB_DSN", "mongodb://prod:27017")
	if got := GetDatabaseDSN(); got != "mongodb://prod:27017" {
		t.Errorf("GetDatabaseDSN() = %q, want %q", got, "mongodb://prod:27017")
	}
}

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

func TestGetFootballDataAPIKey_Empty(t *testing.T) {
	os.Unsetenv("FOOTBALL_DATA_API_KEY")
	if got := GetFootballDataAPIKey(); got != "" {
		t.Errorf("GetFootballDataAPIKey() = %q, want empty", got)
	}
}

func TestGetFootballDataAPIKey_FromEnv(t *testing.T) {
	t.Setenv("FOOTBALL_DATA_API_KEY", "api-key-789")
	if got := GetFootballDataAPIKey(); got != "api-key-789" {
		t.Errorf("GetFootballDataAPIKey() = %q, want %q", got, "api-key-789")
	}
}
