package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}
	defer res.Body.Close()

	var locations []Location
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locations)
	if err != nil {
		log.Fatalf("error decoding response: %v", err)
	}

	logLocations(locations)
}

type Location struct {
	Discovered       bool   `json:"discovered"`
	Id               string `json:"id"`
	Name             string `json:"name"`
	RecommendedLevel int    `json:"recommendedLevel"`
}

func logLocations(locations []Location) {
	for _, l := range locations {
		fmt.Printf("Location: %s, Recommended Character Level: %v\n", l.Name, l.RecommendedLevel)
	}
}
