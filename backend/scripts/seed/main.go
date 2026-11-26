// Command seed populates the matches collection from the frontend's
// worldcup2026.json, upserting by match number.
//
// By default every match is saved as SCHEDULED with no score, as the issue
// asks: results arrive through sync_results or the admin endpoint. With
// -results, the scores in the JSON are loaded and the matches marked FINISHED.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	_ "time/tzdata" // kickoffs are converted with IANA zones; don't rely on the host's zoneinfo.

	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
)

const defaultFile = "../frontend/public/data/worldcup2026.json"

type worldCupFile struct {
	Matches []jsonMatch `json:"matches"`
}

type jsonMatch struct {
	ID      int     `json:"id"`
	Team1   string  `json:"team1"`
	Team2   string  `json:"team2"`
	Goals1  *int    `json:"goals1"`
	Goals2  *int    `json:"goals2"`
	Date    string  `json:"date"`
	Time    string  `json:"time"`
	Stadium string  `json:"stadium"`
	City    string  `json:"city"`
	Group   *string `json:"group"`
	Stage   string  `json:"stage"`
}

var stages = map[string]string{
	"group":   models.StageGroup,
	"round32": models.StageRoundOf32,
	"round16": models.StageRoundOf16,
	"quarter": models.StageQuarterFinals,
	"semi":    models.StageSemiFinals,
	"third":   models.StageThirdPlace,
	"final":   models.StageFinal,
}

// cityTimezones maps each host city to its zone: kickoff times in the JSON are
// local to the stadium.
var cityTimezones = map[string]string{
	"Atlanta":          "America/New_York",
	"Boston":           "America/New_York",
	"Filadélfia":       "America/New_York",
	"Miami":            "America/New_York",
	"Nova York":        "America/New_York",
	"Toronto":          "America/Toronto",
	"Dallas":           "America/Chicago",
	"Houston":          "America/Chicago",
	"Kansas City":      "America/Chicago",
	"Cidade do México": "America/Mexico_City",
	"Guadalajara":      "America/Mexico_City",
	"Monterrey":        "America/Monterrey",
	"Los Angeles":      "America/Los_Angeles",
	"São Francisco":    "America/Los_Angeles",
	"Seattle":          "America/Los_Angeles",
	"Vancouver":        "America/Vancouver",
}

func main() {
	file := flag.String("file", defaultFile, "path to worldcup2026.json")
	withResults := flag.Bool("results", false, "load the JSON scores and mark matches FINISHED")
	flag.Parse()

	data, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("read %s: %v", *file, err)
	}

	matches, err := parseMatches(data, *withResults)
	if err != nil {
		log.Fatalf("parse %s: %v", *file, err)
	}

	mc := config.SetupMongo()
	defer config.TeardownMongo(mc)

	inserted, updated, err := repositories.NewMatchRepository(mc).BulkUpsert(matches)
	if err != nil {
		log.Fatalf("upsert matches: %v", err)
	}

	mode := "scheduled, no scores"
	if *withResults {
		mode = "with results"
	}
	log.Printf("seeded %d matches into %s.%s (%s): %d inserted, %d updated",
		len(matches), config.GetDatabaseName(), config.MatchCollection, mode, inserted, updated)
}

// parseMatches converts the JSON matches into models. Without withResults the
// matches carry no status and no goals, so the upsert creates them as
// SCHEDULED and leaves the status and score of existing matches untouched.
func parseMatches(data []byte, withResults bool) ([]*models.Match, error) {
	var wc worldCupFile
	if err := json.Unmarshal(data, &wc); err != nil {
		return nil, err
	}
	if len(wc.Matches) == 0 {
		return nil, fmt.Errorf("no matches found")
	}

	seen := make(map[int]bool, len(wc.Matches))
	matches := make([]*models.Match, 0, len(wc.Matches))
	for _, jm := range wc.Matches {
		if seen[jm.ID] {
			return nil, fmt.Errorf("match %d: duplicated id", jm.ID)
		}
		seen[jm.ID] = true

		m, err := toMatch(jm, withResults)
		if err != nil {
			return nil, fmt.Errorf("match %d: %w", jm.ID, err)
		}
		matches = append(matches, m)
	}
	return matches, nil
}

func toMatch(jm jsonMatch, withResults bool) (*models.Match, error) {
	if jm.ID <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	stage, ok := stages[jm.Stage]
	if !ok {
		return nil, fmt.Errorf("unknown stage %q", jm.Stage)
	}

	date, err := kickoffUTC(jm.Date, jm.Time, jm.City)
	if err != nil {
		return nil, err
	}

	m := &models.Match{
		MatchNumber: jm.ID,
		Team1:       jm.Team1,
		Team2:       jm.Team2,
		Date:        date,
		Stadium:     jm.Stadium,
		City:        jm.City,
		Stage:       stage,
	}
	if jm.Group != nil {
		m.Group = *jm.Group
	}

	if withResults && jm.Goals1 != nil && jm.Goals2 != nil {
		m.Goals1, m.Goals2 = jm.Goals1, jm.Goals2
		m.Status = models.MatchStatusFinished
	}
	return m, nil
}

// kickoffUTC reads a local date and time at the city's zone and returns it in UTC.
func kickoffUTC(date, clock, city string) (time.Time, error) {
	zone, ok := cityTimezones[city]
	if !ok {
		return time.Time{}, fmt.Errorf("no timezone for city %q", city)
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation("2006-01-02 15:04", date+" "+clock, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid kickoff %q %q: %w", date, clock, err)
	}
	return t.UTC(), nil
}
