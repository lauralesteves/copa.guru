//go:generate mockgen -source=user_repository.go -destination=user_repository_mock.go -package=repositories

package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository interface {
	Get(id bson.ObjectID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetByRefreshToken(token string) (*models.User, error)
	Upsert(user *models.User) (*models.User, error)
	UpdateRefreshToken(userID bson.ObjectID, token string, expiresAt time.Time) error
	InvalidateRefreshToken(userID bson.ObjectID) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *config.MongoContext) UserRepository {
	collection := db.Client.Database(db.Database.Name()).Collection(config.UserCollection)
	return &userRepository{collection: collection}
}

func (r *userRepository) Get(id bson.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"auth.googleId": googleID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByRefreshToken(token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"auth.refreshToken": token}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Upsert(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	user.UpdatedAt = now

	filter := bson.M{"email": user.Email}
	update := bson.M{
		"$set": bson.M{
			"email":            user.Email,
			"name":             user.Name,
			"picture":          user.Picture,
			"auth.lastLoginAt": user.Auth.LastLoginAt,
			"updatedAt":        now,
		},
		"$setOnInsert": bson.M{
			"strategy":      user.Strategy,
			"auth.googleId": user.Auth.GoogleID,
			"auth.password": user.Auth.Password,
			"createdAt":     now,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	if result.UpsertedID != nil {
		user.ID = result.UpsertedID.(bson.ObjectID)
		user.CreatedAt = now
	} else {
		existing, err := r.GetByEmail(user.Email)
		if err != nil {
			return nil, err
		}
		user.ID = existing.ID
		user.CreatedAt = existing.CreatedAt
	}

	return user, nil
}

func (r *userRepository) UpdateRefreshToken(userID bson.ObjectID, token string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"auth.refreshToken":          token,
			"auth.refreshTokenExpiresAt": expiresAt,
			"updatedAt":                  time.Now(),
		}},
	)
	return err
}

func (r *userRepository) InvalidateRefreshToken(userID bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{
			"$unset": bson.M{
				"auth.refreshToken":          "",
				"auth.refreshTokenExpiresAt": "",
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}
