package repositories

import (
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

//go:generate mockgen -source=match_repository.go -destination=match_repository_mock.go -package=repositories

type MatchRepository interface {
	FindByID(id bson.ObjectID) (*models.Match, error)
}
