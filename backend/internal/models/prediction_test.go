package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestPredictionToDTO(t *testing.T) {
	points := 5
	now := time.Now()
	p := &Prediction{
		ID:          bson.NewObjectID(),
		UserID:      bson.NewObjectID(),
		MatchID:     bson.NewObjectID(),
		MatchNumber: 12,
		Goals1:      2,
		Goals2:      1,
		Points:      &points,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("without match", func(t *testing.T) {
		dto := p.ToDTO(nil)
		if dto.ID != p.ID.Hex() || dto.MatchID != p.MatchID.Hex() {
			t.Errorf("ids = %s/%s, want %s/%s", dto.ID, dto.MatchID, p.ID.Hex(), p.MatchID.Hex())
		}
		if dto.MatchNumber != 12 || dto.Goals1 != 2 || dto.Goals2 != 1 {
			t.Errorf("dto = %+v, want match 12, 2-1", dto)
		}
		if dto.Points == nil || *dto.Points != 5 {
			t.Errorf("points = %v, want 5", dto.Points)
		}
		if dto.Match != nil {
			t.Errorf("match = %+v, want nil", dto.Match)
		}
	})

	t.Run("with match", func(t *testing.T) {
		goals1, goals2 := 3, 0
		match := &Match{
			Team1: "Brasil", Team2: "Sérvia", Goals1: &goals1, Goals2: &goals2,
			Date: now, Group: "G", Stage: "GROUP_STAGE", Status: MatchStatusFinished,
		}

		dto := p.ToDTO(match)
		if dto.Match == nil {
			t.Fatal("match = nil, want enriched match")
		}
		want := PredictionMatchDTO{
			Team1: "Brasil", Team2: "Sérvia", Goals1: &goals1, Goals2: &goals2,
			Date: now, Group: "G", Stage: "GROUP_STAGE", Status: MatchStatusFinished,
		}
		if *dto.Match != want {
			t.Errorf("match = %+v, want %+v", *dto.Match, want)
		}
	})
}

func TestPredictionFiltersIsEmpty(t *testing.T) {
	if !(PredictionFilters{}).IsEmpty() {
		t.Error("empty filters reported as not empty")
	}
	if (PredictionFilters{Group: "A"}).IsEmpty() {
		t.Error("group filter reported as empty")
	}
}
