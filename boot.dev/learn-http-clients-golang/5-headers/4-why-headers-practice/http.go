package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
)

type Location struct {
	Discovered       bool   `json:"discovered"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	RecommendedLevel int    `json:"recommendedLevel"`
}

func generateKey() string {
	const characters = "ABCDEF0123456789"
	result := ""
	rand.New(rand.NewSource(0))
	for i := 0; i < 16; i++ {
		result += string(characters[rand.Intn(len(characters))])
	}
	return result
}

func getLocationResponse(apiKey, url string) (Location, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	req.Header.Add("X-API-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	var location Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return Location{}, err
	}

	return location, nil
}

func putLocation(apiKey, url string, location Location) error {
	client := &http.Client{}
	data, err := json.Marshal(location)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("X-API-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
