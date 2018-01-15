package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func fixturesForTeam(team string) ([]Fixture, error) {
	if err := isValidAbbrev(team); err != nil {
		return nil, err
	}

	fullTeamName := Teams[team]
	fixtures, err := getFixtures()
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

func automod(c *gin.Context) {
	team := c.Param("team")

	teamFixtures, err := fixturesForTeam(team)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var response bytes.Buffer

	for _, fixture := range teamFixtures {
		response.WriteString("---\nfirst: \"" +
			fixture.MatchDate.Format("January 1, 2006 15:04 -07") + "\"\n" +
			"sticky: false\n" +
			"distinguish: true\n" +
			"title: \"" + string(fixture.HomeTeam.Name) + " vs " + string(fixture.AwayTeam.Name) + " - Match Thread\"\n" +
			"text: |\n" +
			"  Official match discussion thread\n\n")
	}

	log.Printf("Successful request for %s", team)
	c.String(http.StatusOK, response.String())
}
