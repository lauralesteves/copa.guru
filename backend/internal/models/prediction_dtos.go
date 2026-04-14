package models

import "go.mongodb.org/mongo-driver/v2/bson"

type PredictionDTO struct {
	MatchID     bson.ObjectID `json:"matchId"`
	MatchNumber int           `json:"matchNumber"`
	Goals1      int           `json:"goals1"`
	Goals2      int           `json:"goals2"`
}
