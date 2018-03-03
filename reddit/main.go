package reddit

import (
	"bytes"
	"log"
	"mls-scraper/fixtures"
	"mls-scraper/standings"
	"mls-scraper/teams"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func getTeam(c *gin.Context) string {
	team := c.Param("team")
	if err := teams.AbbrevIsValid(team); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return ""
	}
	return team
}

// AutoMod can be accessed at /reddit/:team/automod
func AutoMod(c *gin.Context) {
	team := getTeam(c)
	if team == "" {
		return
	}

	teamFixtures, err := fixtures.For(team)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var response bytes.Buffer

	for _, fixture := range teamFixtures {
		response.WriteString("---\nfirst: \"" +
			fixture.MatchDate.Format("January 2, 2006 15:04 -07") + "\"\n" +
			"sticky: false\n" +
			"distinguish: true\n" +
			"title: \"" + string(fixture.HomeTeam.Name) + " vs " + string(fixture.AwayTeam.Name) + " - Match Thread\"\n" +
			"text: |\n" +
			"  Official match discussion thread\n\n")
	}

	log.Printf("Successful request for %s", team)
	c.String(http.StatusOK, response.String())
}

// Sidebar can be accessed at /reddit/:team/sidebar
func Sidebar(c *gin.Context) {
	team := getTeam(c)
	if team == "" {
		return
	}

	prevCount, _ := strconv.Atoi(c.DefaultQuery("prevCount", "1"))
	nextCount, _ := strconv.Atoi(c.DefaultQuery("nextCount", "5"))

	teamFixtures, err := fixtures.For(team)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	now := time.Now()
	var prevFixtures []fixtures.Fixture
	var nextFixtures []fixtures.Fixture
	nextIndex := -1
	for i, fixture := range teamFixtures {
		if fixture.MatchDate.Time.Before(now) {
			prevFixtures = append(prevFixtures, fixture)
		} else {
			if nextIndex == -1 {
				nextIndex = i
			}
			nextFixtures = append(nextFixtures, fixture)
			if nextIndex+nextCount-1 == i {
				break
			}
		}
	}
	lenPrev := len(prevFixtures)
	if lenPrev > 0 {
		prevFixtures = prevFixtures[lenPrev-prevCount : lenPrev]
	}
	allFixtures := append(prevFixtures, nextFixtures...)
	var response bytes.Buffer
	response.WriteString("Opponent | Date | Time | Result\n---------|:----:|:----:|-------\n")
	for _, fixture := range allFixtures {
		atV := "v"
		opponent := fixture.AwayTeam.Name
		if opponent == "New England Revolution" {
			atV = "@"
			opponent = fixture.HomeTeam.Name
		}
		response.WriteString(atV + " " + string(opponent) + " | " +
			fixture.MatchDate.Time.Format(" 1.02 | 3:04pm | "))
		if fixture.Finished {
			response.WriteString(strconv.Itoa(fixture.HomeScore) + "-" + strconv.Itoa(fixture.AwayScore))
		}
		response.WriteString("\n")
	}

	c.String(http.StatusOK, response.String())
}

// Standings can be accessed at /reddit/:team/standings
func Standings(c *gin.Context) {
	team := getTeam(c)
	if team == "" {
		return
	}
	conference := teams.ConferenceFor(team)

	table, err := standings.GetFor(conference)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var response bytes.Buffer
	response.WriteString("Club | Pts | PPG | W | L | T | GD\n-----|:---:|:---:|:-:|:-:|:-:|:--\n")
	for _, standing := range table {
		abbrv, _ := teams.AbbrevFor(standing.Name)
		as := ""
		if abbrv == team {
			as = "**"
		}
		cols := []string{
			"[](#" + strings.ToUpper(abbrv) + ") " + standing.Name,
			strconv.Itoa(standing.Points),
			strconv.Itoa(standing.Points / standing.GamesPlayed),
			strconv.Itoa(standing.OverallStats.Wins),
			strconv.Itoa(standing.OverallStats.Losses),
			strconv.Itoa(standing.OverallStats.Draws),
			strconv.Itoa(standing.GoalDifferential),
		}
		response.WriteString(as + strings.Join(cols, as+" | "+as) + as + "\n")
	}

	c.String(http.StatusOK, response.String())
}