package config

import (
	"os"
	"testing"
)

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
