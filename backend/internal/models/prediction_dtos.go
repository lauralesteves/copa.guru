package models

import "time"

type PredictionDTO struct {
	ID          string              `json:"id"`
	MatchID     string              `json:"matchId"`
	MatchNumber int                 `json:"matchNumber"`
	Goals1      int                 `json:"goals1"`
	Goals2      int                 `json:"goals2"`
	Points      *int                `json:"points"`
	Match       *PredictionMatchDTO `json:"match,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

// PredictionMatchDTO is the match data embedded in a PredictionDTO.
type PredictionMatchDTO struct {
	Team1  string      `json:"team1"`
	Team2  string      `json:"team2"`
	Goals1 *int        `json:"goals1"`
	Goals2 *int        `json:"goals2"`
	Date   time.Time   `json:"date"`
	Group  string      `json:"group,omitempty"`
	Stage  string      `json:"stage"`
	Status MatchStatus `json:"status"`
}

type MatchPredictionDTO struct {
	Prediction  *PredictionDTO `json:"prediction"`
	UserName    string         `json:"userName"`
	UserPicture string         `json:"userPicture"`
}

// PredictionFilters narrows ListByUserId by fields of the predicted match.
// Empty fields are ignored.
type PredictionFilters struct {
	Stage string
	Group string
}

func (f PredictionFilters) IsEmpty() bool {
	return f.Stage == "" && f.Group == ""
}

// PredictionSummary aggregates a user's predictions. Unscored predictions
// count towards Total only.
type PredictionSummary struct {
	Total          int `json:"total" bson:"total"`
	ExactScores    int `json:"exactScores" bson:"exactScores"`
	CorrectResults int `json:"correctResults" bson:"correctResults"`
	TotalPoints    int `json:"totalPoints" bson:"totalPoints"`
}

// PredictionDistribution counts the predicted outcomes for a match.
type PredictionDistribution struct {
	MatchID string `json:"matchId" bson:"-"`
	Total   int    `json:"total" bson:"total"`
	HomeWin int    `json:"homeWin" bson:"homeWin"`
	Draw    int    `json:"draw" bson:"draw"`
	AwayWin int    `json:"awayWin" bson:"awayWin"`
}

// PredictionInput get the user input in prediction
type PredictionInput struct {
	MatchNumber int `json:"matchNumber"`
	Goals1      int `json:"goals1"`
	Goals2      int `json:"goals2"`
}

// RejectedPrediction get the reason of rejection
type RejectedPrediction struct {
	Index       int    `json:"index"`
	MatchNumber int    `json:"matchNumber"`
	Reason      string `json:"reason"`
}

// BatchResult
type BatchResult struct {
	Saved    int                  `json:"saved"`
	Rejected []RejectedPrediction `json:"rejected"`
}

type MyPredictionsResponse struct {
	Predictions []*PredictionDTO  `json:"predictions"`
	Summary     PredictionSummary `json:"summary"`
}

type MatchPredictionsResponse struct {
	Match       *MatchDTO             `json:"match"`
	Predictions []*MatchPredictionDTO `json:"predictions"`
}
