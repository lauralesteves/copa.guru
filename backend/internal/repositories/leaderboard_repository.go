//go:generate mockgen -source=leaderboard_repository.go -destination=leaderboard_repository_mock.go -package=repositories

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

type LeaderboardRepository interface {
	List(limit, offset int64) ([]*models.LeaderboardEntry, int64, error)
	GetByUserID(userID bson.ObjectID) (*models.LeaderboardEntry, error)
	ReplaceAll(entries []*models.LeaderboardEntry) error
	Count() (int64, error)
}

type leaderboardRepository struct {
	collection *mongo.Collection
}

func NewLeaderboardRepository(db *config.MongoContext) LeaderboardRepository {
	collection := db.Client.Database(db.Database.Name()).Collection(config.LeaderboardCollection)
	return &leaderboardRepository{collection: collection}
}

// List returns a paginated slice of leaderboard entries sorted by rank ascending,
// along with the total count of entries in the collection.
func (r *leaderboardRepository) List(limit, offset int64) ([]*models.LeaderboardEntry, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "rank", Value: 1}}).
		SetSkip(offset).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var entries []*models.LeaderboardEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// GetByUserID returns a single leaderboard entry for the given user.
// Returns nil if the user has no entry in the leaderboard.
func (r *leaderboardRepository) GetByUserID(userID bson.ObjectID) (*models.LeaderboardEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var entry models.LeaderboardEntry
	err := r.collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&entry)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

// ReplaceAll drops the entire leaderboard collection and inserts the new entries.
// Used by rebuild_leaderboard worker to atomically replace the ranking after scoring.
// Recreates indexes after insertion. If entries is empty, only drops the collection.
func (r *leaderboardRepository) ReplaceAll(entries []*models.LeaderboardEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := r.collection.Drop(ctx); err != nil {
		return err
	}

	if len(entries) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(entries))
	for i, e := range entries {
		e.CreatedAt = now
		e.UpdatedAt = now
		docs[i] = e
	}

	if _, err := r.collection.InsertMany(ctx, docs); err != nil {
		return err
	}

	return r.ensureIndexes(ctx)
}

// Count returns the total number of entries in the leaderboard.
func (r *leaderboardRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{})
}

// ensureIndexes creates indexes for userId (unique), scoring tiebreakers
// (totalPoints desc + exactScores desc), and rank ordering.
func (r *leaderboardRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "totalPoints", Value: -1},
				{Key: "exactScores", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "rank", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
