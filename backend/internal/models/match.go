package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "SCHEDULED"
	MatchStatusLive      MatchStatus = "LIVE"
	MatchStatusFinished  MatchStatus = "FINISHED"
	MatchStatusPostponed MatchStatus = "POSTPONED"
	MatchStatusCancelled MatchStatus = "CANCELLED"
)

type Match struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	ExternalID  string        `json:"externalId" bson:"externalId,omitempty"`
	MatchNumber int           `json:"matchNumber" bson:"matchNumber"`
	Team1       string        `json:"team1" bson:"team1"`
	Team2       string        `json:"team2" bson:"team2"`
	Goals1      *int          `json:"goals1" bson:"goals1,omitempty"`
	Goals2      *int          `json:"goals2" bson:"goals2,omitempty"`
	Date        time.Time     `json:"date" bson:"date"`
	Stadium     string        `json:"stadium" bson:"stadium,omitempty"`
	City        string        `json:"city" bson:"city,omitempty"`
	Group       string        `json:"group" bson:"group,omitempty"`
	Stage       string        `json:"stage" bson:"stage"`
	Status      MatchStatus   `json:"status" bson:"status"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
}
