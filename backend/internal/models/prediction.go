package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Points awarded per prediction. Shared by the scoring service and the
// aggregations that count exact scores and correct results.
const (
	PointsExactScore    = 10
	PointsCorrectResult = 5
)

type Prediction struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      bson.ObjectID `json:"userId" bson:"userId"`
	MatchID     bson.ObjectID `json:"matchId" bson:"matchId"`
	MatchNumber int           `json:"matchNumber" bson:"matchNumber"`
	Goals1      int           `json:"goals1" bson:"goals1"`
	Goals2      int           `json:"goals2" bson:"goals2"`
	// Points is nil until the match is scored.
	Points    *int      `json:"points" bson:"points,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// ToDTO converts the prediction to its response shape. When match is not nil,
// the DTO is enriched with the match data.
func (p *Prediction) ToDTO(match *Match) *PredictionDTO {
	dto := &PredictionDTO{
		ID:          p.ID.Hex(),
		MatchID:     p.MatchID.Hex(),
		MatchNumber: p.MatchNumber,
		Goals1:      p.Goals1,
		Goals2:      p.Goals2,
		Points:      p.Points,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if match != nil {
		dto.Match = &PredictionMatchDTO{
			Team1:  match.Team1,
			Team2:  match.Team2,
			Goals1: match.Goals1,
			Goals2: match.Goals2,
			Date:   match.Date,
			Group:  match.Group,
			Stage:  match.Stage,
			Status: match.Status,
		}
	}

	return dto
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
