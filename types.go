package main

import (
	"regexp"
	"strconv"
	"time"
)

// ByKickoff is a type for sorting Fixtures by kickoff
type ByKickoff []Fixture

func (f ByKickoff) Len() int           { return len(f) }
func (f ByKickoff) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByKickoff) Less(i, j int) bool { return f[i].MatchDate.Time.Before(f[j].MatchDate.Time) }

// kickoffTime needed to parse the info
type kickoffTime struct {
	time.Time
}

type teamName string

// Team contains team data
type Team struct {
	ID   int `json:"id,string"`
	Name teamName
}

// Fixture is another term for match or game
type Fixture struct {
	AwayScore     int `json:"awayScore,string,omitempty"`
	AwayTeam      Team
	Finished      bool
	HomeScore     int `json:"homeScore,string,omitempty"`
	HomeTeam      Team
	ID            int         `json:"id,string"`
	MatchDate     kickoffTime `json:"matchDate,string"`
	StatusInt     int         `json:"status,string,omitempty"`
	StatusOfMatch string      `json:"statusOfMatch,omitempty"`
}

// FotMobData has a bunch of extra stuff
type FotMobData struct {
	League struct {
		Fixtures []Fixture
	}
}

func toInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

func (t *kickoffTime) UnmarshalJSON(data []byte) error {
	dateRegExp := regexp.MustCompile(`(\d{4})(\d{2})(\d{2})(\d{2})(\d{2})`)
	match := dateRegExp.FindStringSubmatch(string(data))
	copenhagenTime, _ := time.LoadLocation("Europe/Copenhagen")
	easternTime, _ := time.LoadLocation("America/New_York")
	date := time.Date(toInt(match[1]), time.Month(toInt(match[2])), toInt(match[3]), toInt(match[4]), toInt(match[5]), 0, 0, copenhagenTime)
	t.Time = date.In(easternTime)
	return nil
}

func (n *teamName) UnmarshalText(data []byte) error {
	switch string(data) {
	case "Atlanta United":
		*n = "Atlanta United FC"
	case "Columbus Crew":
		*n = "Columbus Crew SC"
	case "Minnesota United":
		*n = "Minnesota United FC"
	case "New England Rev.":
		*n = "New England Revolution"
	case "Orlando City":
		*n = "Orlando City SC"
	default:
		if err := isValidTeam(string(data)); err != nil {
			return err
		}
		*n = teamName(string(data))
	}
	return nil
}
