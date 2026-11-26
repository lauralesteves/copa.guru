package models

import "time"

type MatchDTO struct {
	ID          string      `json:"id"`
	ExternalID  string      `json:"externalId,omitempty"`
	MatchNumber int         `json:"matchNumber"`
	Team1       string      `json:"team1"`
	Team2       string      `json:"team2"`
	Goals1      *int        `json:"goals1"`
	Goals2      *int        `json:"goals2"`
	Date        time.Time   `json:"date"`
	Stadium     string      `json:"stadium,omitempty"`
	City        string      `json:"city,omitempty"`
	Group       string      `json:"group,omitempty"`
	Stage       string      `json:"stage"`
	Status      MatchStatus `json:"status"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// MatchFilters narrows MatchRepository.FindAll. Empty fields are ignored.
type MatchFilters struct {
	Stage  string
	Group  string
	Status MatchStatus
}

func (f MatchFilters) IsEmpty() bool {
	return f.Stage == "" && f.Group == "" && f.Status == ""
}
