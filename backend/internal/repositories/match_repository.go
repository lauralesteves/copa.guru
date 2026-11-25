//go:generate mockgen -source=match_repository.go -destination=match_repository_mock.go -package=repositories

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

// ErrMatchNotFound is returned by writes that target a match that does not exist.
var ErrMatchNotFound = errors.New("match not found")

type MatchRepository interface {
	Get(id bson.ObjectID) (*models.Match, error)
	FindAll(filters models.MatchFilters) ([]*models.Match, error)
	FindByMatchNumber(matchNumber int) (*models.Match, error)
	FindByExternalID(externalID string) (*models.Match, error)
	FindScheduledBefore(date time.Time) ([]*models.Match, error)
	FindByStatus(statuses []models.MatchStatus) ([]*models.Match, error)

	UpsertByMatchNumber(match *models.Match) error
	UpdateResult(matchNumber, goals1, goals2 int, status models.MatchStatus) error
	BulkUpsert(matches []*models.Match) (inserted, updated int, err error)
}

type matchRepository struct {
	collection *mongo.Collection
}

func NewMatchRepository(db *config.MongoContext) MatchRepository {
	collection := db.Client.Database(db.Database.Name()).Collection(config.MatchCollection)
	return &matchRepository{collection: collection}
}

var (
	byMatchNumber = bson.D{{Key: "matchNumber", Value: 1}}
	byDate        = bson.D{{Key: "date", Value: 1}, {Key: "matchNumber", Value: 1}}
)

func (r *matchRepository) findOne(filter bson.M) (*models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var match models.Match
	err := r.collection.FindOne(ctx, filter).Decode(&match)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) find(filter bson.M, sort bson.D) ([]*models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	matches := []*models.Match{}
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

// Get returns the match with the given id, or nil when there is none.
func (r *matchRepository) Get(id bson.ObjectID) (*models.Match, error) {
	return r.findOne(bson.M{"_id": id})
}

// FindAll returns the matches ordered by match number, narrowed by filters.
func (r *matchRepository) FindAll(filters models.MatchFilters) ([]*models.Match, error) {
	filter := bson.M{}
	if filters.Stage != "" {
		filter["stage"] = filters.Stage
	}
	if filters.Group != "" {
		filter["group"] = filters.Group
	}
	if filters.Status != "" {
		filter["status"] = filters.Status
	}
	return r.find(filter, byMatchNumber)
}

// FindByMatchNumber returns the match with the given number, or nil when there is none.
func (r *matchRepository) FindByMatchNumber(matchNumber int) (*models.Match, error) {
	return r.findOne(bson.M{"matchNumber": matchNumber})
}

// FindByExternalID returns the match with the given football-data.org id, or
// nil when there is none.
func (r *matchRepository) FindByExternalID(externalID string) (*models.Match, error) {
	if externalID == "" {
		return nil, nil
	}
	return r.findOne(bson.M{"externalId": externalID})
}

// FindScheduledBefore returns the matches still SCHEDULED whose kickoff is
// before date, that is, matches that should already have started.
func (r *matchRepository) FindScheduledBefore(date time.Time) ([]*models.Match, error) {
	return r.find(bson.M{
		"status": models.MatchStatusScheduled,
		"date":   bson.M{"$lt": date},
	}, byDate)
}

// FindByStatus returns the matches in any of the given statuses, ordered by date.
func (r *matchRepository) FindByStatus(statuses []models.MatchStatus) ([]*models.Match, error) {
	if len(statuses) == 0 {
		return []*models.Match{}, nil
	}
	return r.find(bson.M{"status": bson.M{"$in": statuses}}, byDate)
}

// upsertByMatchNumber builds the filter and update of an upsert keyed by
// matchNumber. Empty external ids and nil goals are left untouched, so
// re-running a seed never erases a synced id or a result, and the sparse
// unique index on externalId never sees duplicated empty strings.
func upsertByMatchNumber(match *models.Match, now time.Time) (bson.M, bson.M) {
	set := bson.M{
		"team1":     match.Team1,
		"team2":     match.Team2,
		"date":      match.Date,
		"stadium":   match.Stadium,
		"city":      match.City,
		"group":     match.Group,
		"stage":     match.Stage,
		"status":    match.Status,
		"updatedAt": now,
	}
	if match.ExternalID != "" {
		set["externalId"] = match.ExternalID
	}
	if match.Goals1 != nil {
		set["goals1"] = *match.Goals1
	}
	if match.Goals2 != nil {
		set["goals2"] = *match.Goals2
	}

	return bson.M{"matchNumber": match.MatchNumber}, bson.M{"$set": set}
}

// UpsertByMatchNumber creates or updates a match keyed by its match number.
func (r *matchRepository) UpsertByMatchNumber(match *models.Match) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter, update := upsertByMatchNumber(match, time.Now())
	_, err := r.collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	return err
}

// UpdateResult sets the score and status of a match. It returns
// ErrMatchNotFound when no match has the given number.
func (r *matchRepository) UpdateResult(matchNumber, goals1, goals2 int, status models.MatchStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(ctx,
		bson.M{"matchNumber": matchNumber},
		bson.M{"$set": bson.M{
			"goals1":    goals1,
			"goals2":    goals2,
			"status":    status,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMatchNotFound
	}
	return nil
}

// BulkUpsert creates or updates several matches in one round trip, keyed by
// match number. It returns how many matches were inserted and how many
// existing ones were updated.
func (r *matchRepository) BulkUpsert(matches []*models.Match) (int, int, error) {
	if len(matches) == 0 {
		return 0, 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	operations := make([]mongo.WriteModel, 0, len(matches))
	for _, m := range matches {
		filter, update := upsertByMatchNumber(m, now)
		operations = append(operations, mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true))
	}

	result, err := r.collection.BulkWrite(ctx, operations, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return 0, 0, err
	}
	return int(result.UpsertedCount), int(result.MatchedCount), nil
}
