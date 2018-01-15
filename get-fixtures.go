package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"sort"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const fotMobURL = "https://www.fotmob.com/leagues/130/matches/"

func getFixtures() ([]Fixture, error) {
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
