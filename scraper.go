package main

import (
	"encoding/json"
	"fmt"
)

func scrapeFotMob() {
	fixtures, err := getFixtures()
	if err != nil {
		panic(err)
	}

	jsonStr, err := json.Marshal(fixtures)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonStr))
}

func main() {
	scrapeFotMob()
}
