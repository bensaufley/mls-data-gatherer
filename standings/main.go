package standings

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Shield gives Supporters Shield Standings
func Shield(c *gin.Context) {
	standings, err := GetShield()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, gin.H{
		"standings": standings,
	})
}

// Conference gives standings by Conference
func Conference(c *gin.Context) {
	conference := strings.ToLower(c.Param("conference"))
	if conference != "east" && conference != "west" {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid Conference"))
	}
	standings, err := GetFor(conference)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, gin.H{
		"standings": standings,
	})
}
