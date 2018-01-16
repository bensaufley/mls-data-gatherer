package teams

// Conferences maps team abbreviations to conferences
var Conferences = map[string][]string{
	"east": []string{"atl", "chi", "clb", "dc", "mtl", "ner", "nyc", "nyrb", "orl", "phi", "tfc"},
	"west": []string{"dal", "hou", "col", "lag", "lafc", "min", "por", "rsl", "sj", "sea", "kc", "van"},
}

// ConferenceFor returns the conference associated with a team
func ConferenceFor(abbrv string) string {
	for conference, teams := range Conferences {
		for _, team := range teams {
			if team == abbrv {
				return conference
			}
		}
	}
	return ""
}

// Name is the name of the team
type Name string

// Team contains team data
type Team struct {
	ID   int `json:"id,string"`
	Name Name
}

// Teams is the canonical list of team names and abbreviations
var Teams = map[string]string{"atl": "Atlanta United FC",
	"chi":  "Chicago Fire",
	"col":  "Colorado Rapids",
	"clb":  "Columbus Crew SC",
	"dc":   "D.C. United",
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
	"van":  "Vancouver Whitecaps FC",
}

// UnmarshalText Coerces team Name to the right one
func (n *Name) UnmarshalText(data []byte) error {
	switch string(data) {
	case "Atlanta United":
		*n = "Atlanta United FC"
	case "Columbus Crew":
		*n = "Columbus Crew SC"
	case "DC United":
		*n = "D.C. United"
	case "Minnesota United":
		*n = "Minnesota United FC"
	case "New England Rev.":
		*n = "New England Revolution"
	case "Orlando City":
		*n = "Orlando City SC"
	case "Vancouver Whitecaps":
		*n = "Vancouver Whitecaps FC"
	default:
		if err := NameIsValid(string(data)); err != nil {
			return err
		}
		*n = Name(string(data))
	}
	return nil
}
