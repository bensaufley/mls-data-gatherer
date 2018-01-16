package main

import (
	"log"
	"mls-scraper/fixtures"
	"mls-scraper/reddit"
	"mls-scraper/standings"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/reddit/:team/automod", reddit.AutoMod)
	router.GET("/reddit/:team/sidebar", reddit.Sidebar)
	router.GET("/reddit/:team/standings", reddit.Standings)
	router.GET("/standings/shield", standings.Shield)
	router.GET("/standings/conference/:conference", standings.Conference)
	router.GET("/", func(c *gin.Context) {
		fixtures, err := fixtures.Get()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"fixtures": fixtures,
		})
	})
	router.Run(":" + port)
}
