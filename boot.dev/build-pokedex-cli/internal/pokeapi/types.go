package pokeapi

import (
	"math/rand"
)

type resLocationArea struct {
	// This response includes a LOT of fields that we don't use so we only unmarshal the ones we need
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type resLocationAreaPage struct {
	// Fields have to be exported (public) in order to marshal/unmarshal them
	// So we can't just use lowercase field names that match the json fields
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type resPokemon struct {
	// Again, including only the fields we're going to use
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type Pokemon resPokemon

func (p Pokemon) TryCatch() bool {
	// assumes a standard pokeball...
	// Usually base experience goes from 50 ~ 350
	return (rand.Intn(400) - 50) > p.BaseExperience
}
