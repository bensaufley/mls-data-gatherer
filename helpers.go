package main

import "errors"

// Teams is the canonical list of team names and abbreviations
var Teams = map[string]string{"atl": "Atlanta United FC",
	"chi":  "Chicago Fire",
	"col":  "Colorado Rapids",
	"clb":  "Columbus Crew SC",
	"dc":   "DC United",
	"dal":  "FC Dallas",
	"hou":  "Houston Dynamo",
	"lag":  "LA Galaxy",
	"lafc": "Los Angeles FC",
	"min":  "Minnesota United FC",
	"mtl":  "Montreal Impact",
	"ner":  "New England Revolution",
	"nyc":  "New York City FC",
	"nyrb": "New York Red Bulls",
	"orl":  "Orlando City SC",
	"phi":  "Philadelphia Union",
	"por":  "Portland Timbers",
	"rsl":  "Real Salt Lake",
	"sj":   "San Jose Earthquakes",
	"sea":  "Seattle Sounders FC",
	"kc":   "Sporting Kansas City",
	"tfc":  "Toronto FC",
	"van":  "Vancouver Whitecaps",
}

func isValidAbbrev(team string) error {
	for name := range Teams {
		if name == team {
			return nil
		}
	}
	return errors.New("No such team abbreviation " + team)
}

func isValidTeam(team string) error {
	for _, name := range Teams {
		if name == team {
			return nil
		}
	}
	return errors.New("No such team name " + team)
}
