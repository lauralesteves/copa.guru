package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: create-local-user email:password")
	}

	parts := strings.SplitN(os.Args[1], ":", 2)
	if len(parts) != 2 {
		log.Fatal("format must be email:password")
	}
	email, password := parts[0], parts[1]

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(config.GetDatabaseDSN()))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.GetDatabaseName()).Collection("users")
	now := time.Now()

	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"name":      strings.Split(email, "@")[0],
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"strategy":  "email",
			"email":     email,
			"password":  string(hash),
			"createdAt": now,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf("failed to upsert user: %v", err)
	}

	if result.UpsertedID != nil {
		fmt.Printf("user created: %s\n", email)
	} else {
		fmt.Printf("user already exists: %s\n", email)
	}
}
