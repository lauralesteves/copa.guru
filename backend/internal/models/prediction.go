package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Prediction struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      bson.ObjectID `json:"userId" bson:"userId"`
	MatchID     bson.ObjectID `json:"matchId" bson:"matchId"`
	MatchNumber int           `json:"matchNumber" bson:"matchNumber"`
	Goals1      int           `json:"goals1" bson:"goals1"`
	Goals2      int           `json:"goals2" bson:"goals2"`
	Points      *int          `json:"points" bson:"points,omitempty"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type PointsUpdate struct {
	PredictionID bson.ObjectID
	Points       int
}

type ScoreResult struct {
	Total          int `json:"total"`
	ExactScores    int `json:"exactScores"`
	CorrectResults int `json:"correctResults"`
}

func (p *Prediction) ToDTO() *PredictionDTO {
	return &PredictionDTO{
		MatchID:     p.MatchID,
		MatchNumber: p.MatchNumber,
		Goals1:      p.Goals1,
		Goals2:      p.Goals2,
	}
}
