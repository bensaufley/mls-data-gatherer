package fixtures

import (
	"encoding/json"
	"errors"
	"mls-data-gatherer/teams"
	"net/http"
	"net/url"
	"regexp"
	"sort"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const fotMobURL = "https://www.fotmob.com/leagues/130/matches/"

// For returns Fixtures for a team
func For(team string) ([]Fixture, error) {
	if err := teams.AbbrevIsValid(team); err != nil {
		return nil, err
	}

	fullTeamName := teams.Teams[team]
	fixtures, err := Get()
	if err != nil {
		return nil, err
	}

	var teamFixtures []Fixture

	for _, fixture := range fixtures {
		if string(fixture.AwayTeam.Name) == fullTeamName || string(fixture.HomeTeam.Name) == fullTeamName {
			teamFixtures = append(teamFixtures, fixture)
		}
	}

	if len(teamFixtures) == 0 {
		return nil, errors.New("No fixtures found")
	}

	return teamFixtures, nil
}

// Get gets array of all Fixtures
func Get() ([]Fixture, error) {
	scriptRegExp := regexp.MustCompile(`[\s\S]*window\.__INITIAL_STATE__ *= *"(?P<Data>[\S][\s\S]+[^\\])";[\s\S]*`)

	resp, err := http.Get(fotMobURL)
	if err != nil {
		return nil, err
	}
	root, _ := html.Parse(resp.Body)
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Script {
			return scriptRegExp.MatchString(scrape.Text(n))
		}
		return false
	}
	script, found := scrape.Find(root, matcher)
	if found == false {
		return nil, errors.New("Could not find data")
	}
	scriptContent := scriptRegExp.ReplaceAllString(scrape.Text(script), `${Data}`)

	decoded, _ := url.QueryUnescape(scriptContent)
	var data FotMobData
	if err = json.Unmarshal([]byte(decoded), &data); err != nil {
		return nil, err
	}
	fixtures := data.League.Fixtures
	sort.Sort(ByKickoff(fixtures))
	return fixtures, nil
}
