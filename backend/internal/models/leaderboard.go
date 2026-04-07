package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LeaderboardEntry struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      bson.ObjectID `json:"userId" bson:"userId"`
	UserName    string        `json:"userName" bson:"userName"`
	UserPicture string        `json:"userPicture" bson:"userPicture,omitempty"`

	TotalPoints      int `json:"totalPoints" bson:"totalPoints"`
	ExactScores      int `json:"exactScores" bson:"exactScores"`
	CorrectResults   int `json:"correctResults" bson:"correctResults"`
	TotalPredictions int `json:"totalPredictions" bson:"totalPredictions"`
	Rank             int `json:"rank" bson:"rank"`

	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (e *LeaderboardEntry) ToDTO() *LeaderboardEntryDTO {
	return &LeaderboardEntryDTO{
		UserID:           e.UserID.Hex(),
		UserName:         e.UserName,
		UserPicture:      e.UserPicture,
		TotalPoints:      e.TotalPoints,
		ExactScores:      e.ExactScores,
		CorrectResults:   e.CorrectResults,
		TotalPredictions: e.TotalPredictions,
		Rank:             e.Rank,
	}
}
