package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const fotMobURL = "https://www.fotmob.com/leagues/130/matches/"

func scrape() {
	resp, err := http.Get(fotMobURL)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n.Parent.Parent, "class") == "athing"
		}
		return false
	}
	articles := scrape.FindAll(root, matcher)
	for i, article := range articles {
		fmt.Printf("%2d %s (%s)\n", i, scrape.Text(article), scrape.Attr(article, "href"))
	}
}

func main() {
	scrape()
}
