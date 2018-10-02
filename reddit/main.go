package reddit

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bensaufley/mls-data-gatherer/fixtures"
	"github.com/bensaufley/mls-data-gatherer/standings"
	"github.com/bensaufley/mls-data-gatherer/teams"

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

func getOffset(c *gin.Context) time.Duration {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	return time.Duration(offset)
}

// AutoMod can be accessed at /reddit/:team/automod
func AutoMod(c *gin.Context) {
	team := getTeam(c)
	if team == "" {
		return
	}
	offset := getOffset(c)

	teamFixtures, err := fixtures.For(team)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var response bytes.Buffer

	for _, fixture := range teamFixtures {
		response.WriteString("---\n    first: \"" +
			(fixture.MatchDate.Add(-60 * offset * time.Minute)).Format("January 2, 2006 15:04 -07") + "\"\n" +
			"    sticky: false\n" +
			"    distinguish: true\n" +
			"    title: \"" + string(fixture.HomeTeam.Name) + " vs " + string(fixture.AwayTeam.Name) + " - Match Thread\"\n" +
			"    text: |\n" +
			"      Official match discussion thread\n\n")
	}

	log.Printf("Successful request for %s", team)
	c.String(http.StatusOK, response.String())
}

// Schedule can be accessed at /reddit/:team/schedule
func Schedule(c *gin.Context) {
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
			result := "D"
			if atV == "@" {
				if *fixture.HomeScore > *fixture.AwayScore {
					result = "L"
				} else if *fixture.HomeScore < *fixture.AwayScore {
					result = "W"
				}
			} else {
				if *fixture.HomeScore > *fixture.AwayScore {
					result = "W"
				} else if *fixture.HomeScore < *fixture.AwayScore {
					result = "L"
				}
			}
			response.WriteString(result + " " + strconv.Itoa(*fixture.HomeScore) + "-" + strconv.Itoa(*fixture.AwayScore))
		}
		response.WriteString("\n")
	}

	c.String(http.StatusOK, response.String())
}

// Standings can be accessed at /reddit/:team/standings
func Standings(c *gin.Context) {
	requestedTeam := getTeam(c)
	if requestedTeam == "" {
		return
	}
	conference := teams.ConferenceFor(requestedTeam)

	table, err := standings.GetFor(conference)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var response bytes.Buffer
	response.WriteString("Club | Pts | PPG | W | L | T | GD\n-----|:---:|:---:|:-:|:-:|:-:|:--\n")
	for _, team := range table {
		abbrv, _ := teams.AbbrevFor(team.Name)
		as := ""
		if abbrv == requestedTeam {
			as = "**"
		}
		gpg := 0.0
		if team.GamesPlayed > 0 {
			gpg = float64(team.Points) / float64(team.GamesPlayed)
		}
		cols := []string{
			"[](#" + strings.ToUpper(abbrv) + ") " + team.Name,
			strconv.Itoa(team.Points),
			strconv.FormatFloat(gpg, 'f', 2, 64),
			strconv.Itoa(team.OverallStats.Wins),
			strconv.Itoa(team.OverallStats.Losses),
			strconv.Itoa(team.OverallStats.Draws),
			strconv.Itoa(team.GoalDifferential),
		}
		response.WriteString(as + strings.Join(cols, as+" | "+as) + as + "\n")
	}

	c.String(http.StatusOK, response.String())
}
