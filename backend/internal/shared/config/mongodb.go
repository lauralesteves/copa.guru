package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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

type MongoContext struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func ConnectOnMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(GetDatabaseDSN()))
	if err != nil {
		log.Panicf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Panicf("failed to ping MongoDB: %v", err)
	}

	return client
}

func SetupMongo() *MongoContext {
	client := ConnectOnMongo()
	db := client.Database(GetDatabaseName())
	return &MongoContext{Client: client, Database: db}
}

func TeardownMongo(mc *MongoContext) {
	if mc == nil || mc.Client == nil {
		return
	}
	if err := mc.Client.Disconnect(context.Background()); err != nil {
		log.Printf("failed to disconnect from MongoDB: %v", err)
	}
}
