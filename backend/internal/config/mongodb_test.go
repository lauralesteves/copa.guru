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
