package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lauralesteves/copa-guru-backend/internal/models"
)

// The real file, relative to this package.
const realFile = "../../" + defaultFile

func loadReal(t *testing.T, withResults bool) []*models.Match {
	t.Helper()
	data, err := os.ReadFile(realFile)
	if err != nil {
		t.Fatalf("read %s: %v", realFile, err)
	}
	matches, err := parseMatches(data, withResults)
	if err != nil {
		t.Fatalf("parseMatches() error: %v", err)
	}
	return matches
}

func TestParseMatches_RealFile(t *testing.T) {
	matches := loadReal(t, false)
	if len(matches) != 104 {
		t.Fatalf("len = %d, want 104", len(matches))
	}

	for _, m := range matches {
		if m.Status != "" || m.Goals1 != nil || m.Goals2 != nil {
			t.Errorf("match %d = status %q, goals %v-%v, want no status and no score", m.MatchNumber, m.Status, m.Goals1, m.Goals2)
		}
		if m.Date.Location() != time.UTC {
			t.Errorf("match %d date in %s, want UTC", m.MatchNumber, m.Date.Location())
		}
		if (m.Stage == models.StageGroup) != (m.Group != "") {
			t.Errorf("match %d: stage %s with group %q", m.MatchNumber, m.Stage, m.Group)
		}
	}

	// Opener at Estadio Azteca, 13:00 in Mexico City (UTC-6).
	if want := time.Date(2026, 6, 11, 19, 0, 0, 0, time.UTC); !matches[0].Date.Equal(want) {
		t.Errorf("opener date = %s, want %s", matches[0].Date, want)
	}
	// Final at MetLife Stadium, 15:00 in New York (UTC-4 in July).
	final := matches[103]
	if final.Stage != models.StageFinal {
		t.Fatalf("match 104 stage = %s, want FINAL", final.Stage)
	}
	if want := time.Date(2026, 7, 19, 19, 0, 0, 0, time.UTC); !final.Date.Equal(want) {
		t.Errorf("final date = %s, want %s", final.Date, want)
	}
}

func TestParseMatches_WithResults(t *testing.T) {
	for _, m := range loadReal(t, true) {
		if m.Status != models.MatchStatusFinished || m.Goals1 == nil || m.Goals2 == nil {
			t.Errorf("match %d = status %q, goals %v-%v, want FINISHED with score", m.MatchNumber, m.Status, m.Goals1, m.Goals2)
		}
	}
}

func TestParseMatches_ResultsWithoutScore(t *testing.T) {
	data := `{"matches":[{"id":1,"team1":"MEX","team2":"RSA","goals1":null,"goals2":null,
		"date":"2026-06-11","time":"13:00","city":"Cidade do México","group":"A","stage":"group"}]}`

	matches, err := parseMatches([]byte(data), true)
	if err != nil {
		t.Fatalf("parseMatches() error: %v", err)
	}
	if matches[0].Status != "" {
		t.Errorf("status = %q, want none for a match without score", matches[0].Status)
	}
}

func TestParseMatches_Errors(t *testing.T) {
	valid := `"team1":"MEX","team2":"RSA","date":"2026-06-11","time":"13:00","city":"Cidade do México","stage":"group"`

	tests := []struct {
		name, data, want string
	}{
		{"invalid json", `{`, "unexpected end"},
		{"no matches", `{"matches":[]}`, "no matches"},
		{"invalid id", `{"matches":[{"id":0,` + valid + `}]}`, "invalid id"},
		{"duplicated id", `{"matches":[{"id":1,` + valid + `},{"id":1,` + valid + `}]}`, "duplicated id"},
		{"unknown stage", `{"matches":[{"id":1,"stage":"playoff","city":"Dallas","date":"2026-06-11","time":"13:00"}]}`, "unknown stage"},
		{"unknown city", `{"matches":[{"id":1,"stage":"group","city":"Lisboa","date":"2026-06-11","time":"13:00"}]}`, "no timezone"},
		{"invalid time", `{"matches":[{"id":1,"stage":"group","city":"Dallas","date":"2026-06-11","time":"1pm"}]}`, "invalid kickoff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseMatches([]byte(tt.data), false)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("parseMatches() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}
