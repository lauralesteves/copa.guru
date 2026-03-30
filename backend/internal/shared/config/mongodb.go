package config

import "os"

const (
	defaultDatabaseName = "copa-guru"
	defaultDSN          = "mongodb://localhost:27017"
)

func GetDatabaseName() string {
	if v := os.Getenv("MONGODB_DATABASE_NAME"); v != "" {
		return v
	}
	return defaultDatabaseName
}

func GetDatabaseDSN() string {
	if v := os.Getenv("MONGODB_DSN"); v != "" {
		return v
	}
	return defaultDSN
}
