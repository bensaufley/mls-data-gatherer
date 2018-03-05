package standings

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func extractStandingFromNode(row *html.Node, headers []string) Standing {
	var standing Standing
	teamRegExp := regexp.MustCompile(`^(?:([syx]) - )?(?:[A-Z]+) (.+)$`)
	cols := scrape.FindAll(row, scrape.ByTag(atom.Td))
	for j, col := range cols {
		switch headers[j] {
		case "#":
			standing.Place = intFromScrape(col)
		case "Club":
			content := scrape.Text(col)
			pieces := teamRegExp.FindStringSubmatch(content)
			switch pieces[1] {
			case "s":
				standing.ShieldWinner = true
				standing.ConferenceWinner = true
				standing.Clinched = true
			case "y":
				standing.ConferenceWinner = true
				standing.Clinched = true
			case "x":
				standing.Clinched = true
			}
			standing.Name = pieces[2]
		case "PTS":
			standing.Points = intFromScrape(col)
		case "GP":
			standing.GamesPlayed = intFromScrape(col)
		case "W":
			standing.OverallStats.Wins = intFromScrape(col)
		case "L":
			standing.OverallStats.Losses = intFromScrape(col)
		case "T":
			standing.OverallStats.Draws = intFromScrape(col)
		case "GF":
			standing.GoalsFor = intFromScrape(col)
		case "GA":
			standing.GoalsAgainst = intFromScrape(col)
		case "GD":
			standing.GoalDifferential = intFromScrape(col)
		case "W-L-T 0":
			standing.HomeStats = getWltFromNode(col)
		case "W-L-T 1":
			standing.AwayStats = getWltFromNode(col)
		}
	}
	return standing
}

func extractStandingsFromTable(table *html.Node) Standings {
	tbody, _ := scrape.Find(table, scrape.ByTag(atom.Tbody))
	trs := scrape.FindAll(tbody, scrape.ByTag(atom.Tr))
	headerRow, standingRows := trs[0], trs[1:]
	headerCols := scrape.FindAll(headerRow, scrape.ByTag(atom.Td))
	var headers []string
	wlt := 0
	for _, col := range headerCols {
		text := scrape.Text(col)
		if text == "W-L-T" {
			headers = append(headers, "W-L-T "+strconv.Itoa(wlt))
			wlt++
		} else {
			headers = append(headers, text)
		}
	}
	var standings Standings
	for _, row := range standingRows {
		standings = append(standings, extractStandingFromNode(row, headers))
	}
	return standings
}

func getWltFromNode(node *html.Node) Stats {
	wltRegExp := regexp.MustCompile(`^\D*(\d+)\D+(\d+)\D+(\d+)\D*$`)
	record := scrape.Text(node)
	pieces := wltRegExp.FindStringSubmatch(record)
	var stats Stats
	stats.Wins, _ = strconv.Atoi(pieces[1])
	stats.Losses, _ = strconv.Atoi(pieces[2])
	stats.Draws, _ = strconv.Atoi(pieces[3])
	return stats
}

func intFromScrape(node *html.Node) int {
	i, err := strconv.Atoi(scrape.Text(node))
	if err != nil {
		return 0
	}
	return i
}

// GetShield gets shield standings
func GetShield() (Standings, error) {
	resp, err := http.Get("https://www.mlssoccer.com/standings/supporters-shield")
	if err != nil {
		return nil, err
	}
	root, _ := html.Parse(resp.Body)
	table, ok := scrape.Find(root, scrape.ByClass("standings_table"))
	if ok == false {
		return nil, errors.New("No standings table found")
	}
	standings := extractStandingsFromTable(table)
	return standings, nil
}

// GetFor gets standings for conference
func GetFor(conference string) (Standings, error) {
	resp, err := http.Get("https://www.mlssoccer.com/standings")
	if err != nil {
		return nil, err
	}
	root, _ := html.Parse(resp.Body)
	tables := scrape.FindAll(root, scrape.ByClass("standings_table"))
	var table *html.Node
	confName := conference + "ern conference"
	for _, t := range tables {
		h2, ok := scrape.FindPrevSibling(t, scrape.ByTag(atom.H2))
		if ok && strings.ToLower(scrape.Text(h2)) == confName {
			table = t
			break
		}
	}
	if table == nil {
		return nil, errors.New("No standings table found")
	}
	standings := extractStandingsFromTable(table)
	return standings, nil
}
