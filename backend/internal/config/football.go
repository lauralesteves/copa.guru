package config

import "os"

func GetFootballDataAPIKey() string {
	return os.Getenv("FOOTBALL_DATA_API_KEY")
}
