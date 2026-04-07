package models

type LeaderboardEntryDTO struct {
	UserID           string `json:"userId"`
	UserName         string `json:"userName"`
	UserPicture      string `json:"userPicture"`
	TotalPoints      int    `json:"totalPoints"`
	ExactScores      int    `json:"exactScores"`
	CorrectResults   int    `json:"correctResults"`
	TotalPredictions int    `json:"totalPredictions"`
	Rank             int    `json:"rank"`
}
