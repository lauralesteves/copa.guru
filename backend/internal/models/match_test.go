package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestMatchToDTO(t *testing.T) {
	goals1, goals2 := 2, 0
	now := time.Now()
	m := &Match{
		ID:          bson.NewObjectID(),
		ExternalID:  "537327",
		MatchNumber: 1,
		Team1:       "México",
		Team2:       "África do Sul",
		Goals1:      &goals1,
		Goals2:      &goals2,
		Date:        now,
		Stadium:     "Estadio Azteca",
		City:        "Cidade do México",
		Group:       "A",
		Stage:       "GROUP_STAGE",
		Status:      MatchStatusFinished,
		UpdatedAt:   now,
	}

	want := MatchDTO{
		ID:          m.ID.Hex(),
		ExternalID:  "537327",
		MatchNumber: 1,
		Team1:       "México",
		Team2:       "África do Sul",
		Goals1:      &goals1,
		Goals2:      &goals2,
		Date:        now,
		Stadium:     "Estadio Azteca",
		City:        "Cidade do México",
		Group:       "A",
		Stage:       "GROUP_STAGE",
		Status:      MatchStatusFinished,
		UpdatedAt:   now,
	}
	if got := m.ToDTO(); *got != want {
		t.Errorf("ToDTO() = %+v, want %+v", *got, want)
	}
}

func TestMatchToDTO_Unplayed(t *testing.T) {
	dto := (&Match{Status: MatchStatusScheduled}).ToDTO()
	if dto.Goals1 != nil || dto.Goals2 != nil {
		t.Errorf("goals = %v-%v, want nil until the match is played", dto.Goals1, dto.Goals2)
	}
}

func TestMatchStatusIsValid(t *testing.T) {
	for _, s := range []MatchStatus{
		MatchStatusScheduled, MatchStatusLive, MatchStatusFinished, MatchStatusPostponed, MatchStatusCancelled,
	} {
		if !s.IsValid() {
			t.Errorf("%s reported as invalid", s)
		}
	}
	for _, s := range []MatchStatus{"", "finished", "IN_PLAY"} {
		if s.IsValid() {
			t.Errorf("%q reported as valid", s)
		}
	}
}

func TestMatchFiltersIsEmpty(t *testing.T) {
	if !(MatchFilters{}).IsEmpty() {
		t.Error("empty filters reported as not empty")
	}
	if (MatchFilters{Status: MatchStatusLive}).IsEmpty() {
		t.Error("status filter reported as empty")
	}
}
