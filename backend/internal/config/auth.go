package config

import "os"

const (
	defaultJWTSecret = "copa-guru-local-jwt-secret-key"
)

func GetJWTSecret() string {
	if v := os.Getenv("AUTH_JWT_SECRET"); v != "" {
		return v
	}
	return defaultJWTSecret
}

func GetGoogleClientID() string {
	return os.Getenv("GOOGLE_CLIENT_ID")
}

func GetGoogleClientSecret() string {
	return os.Getenv("GOOGLE_CLIENT_SECRET")
}
