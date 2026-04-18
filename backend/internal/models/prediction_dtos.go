package models

import "go.mongodb.org/mongo-driver/v2/bson"

type PredictionDTO struct {
	MatchID     string `json:"matchId"`
	MatchNumber int    `json:"matchNumber"`
	Goals1      int    `json:"goals1"`
	Goals2      int    `json:"goals2"`
}

type PredictionFilters struct {
	Stage *string `json:"stage,omitempty"`
	Group *string `json:"group,omitempty"`
}

type PredictionDistribution struct {
	HomeWin int    `json:"homeWin"`
	Draw    int    `json:"draw"`
	AwayWin int    `json:"awayWin"`
	MatchID string `json:"matchId"`
	Entries int    `json:"entries"`
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

type UserSummary struct {
	UserID         string `json:"userId"`
	TotalPoints    int    `json:"totalPoints"`
	ExactScores    int    `json:"exactScores"`
	CorrectResults int    `json:"correctResults"`
}
