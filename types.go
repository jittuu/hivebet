package hivebet

import (
	"time"
)

// Event represents a soccer match event.
type Event struct {
	ID        int64 `datastore:"-" goon:"id"`
	League    League
	Season    string
	StartTime time.Time
	Home      string
	Away      string
	HGoal     int   `datastore:",noindex"`
	HGoals    []int `datastore:",noindex"`
	AGoals    []int `datastore:",noindex"`
	AGoal     int   `datastore:",noindex"`
	AvgOdds   MatchOdds
	MaxOdds   MatchOdds
	AvgAHOdds AHOdds
	MaxAHOdds AHOdds
	AvgOUOdds OUOdds
	MaxOUOdds OUOdds
}

// MatchOdds represents 1x2 odds
type MatchOdds struct {
	Home float64 `datastore:",noindex"`
	Draw float64 `datastore:",noindex"`
	Away float64 `datastore:",noindex"`
}

// AHOdds represents Asia Handicap odds
type AHOdds struct {
	Home     float64 `datastore:",noindex"`
	Away     float64 `datastore:",noindex"`
	Handicap float64 `datastore:",noindex"`
}

// IsHomeFavorite returns true if home is favorite team
func (ah *AHOdds) IsHomeFavorite() bool {
	if ah.Handicap < 0 {
		return true
	}

	if ah.Handicap == 0 {
		return ah.Home < ah.Away
	}

	return false
}

// OUOdds represents Over Under odds
type OUOdds struct {
	Over     float64 `datastore:",noindex"`
	Under    float64 `datastore:",noindex"`
	Handicap float64 `datastore:",noindex"`
}

// League is soccer league. E.g. EPL, Serie-A
type League string

// EPL is england premier league
const EPL League = "epl"

// SerieA is italy league
const SerieA League = "serie-a"

// Bundesliga is german league
const Bundesliga League = "bundesliga"

// LaLiga is spain league
const LaLiga = "la-liga"
