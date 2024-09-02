package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mosamadeeb/pokedexcli/internal/pokecache"
)

const (
	apiBase         = "https://pokeapi.co/api/v2/"
	apiLocationArea = apiBase + "location-area/"
	pageSize        = 20

	cacheReapDuration = 5 * time.Minute
)

var cache pokecache.Cache

type PageConfig struct {
	Next     *string
	Previous *string
}

func init() {
	cache = pokecache.NewCache(cacheReapDuration)
}

func apiGet(url string) ([]byte, error) {
	body, ok := cache.Get(url)

	if ok {
		return body, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get request error: %w", err)
	}

	// Important: defer closing the body reader *after* checking for the error
	// If the error is nil, the body will be already closed
	// If the error is not nil, closing the body will allow us to reuse the connection
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected api response status: %s", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	cache.Add(url, body)
	return body, nil
}

func FetchLocationArea(page *PageConfig, fetchPrevPage bool) ([]string, error) {
	var url string
	if fetchPrevPage {
		if page.Previous == nil {
			return nil, errors.New("no previous page available")
		}

		url = *page.Previous
	} else {
		if page.Next == nil {
			// This means we need to fetch the first page
			url = apiLocationArea
		} else {
			url = *page.Next
		}
	}

	body, err := apiGet(url)
	if err != nil {
		return nil, err
	}

	resData := resLocationAreaPage{}
	err = json.Unmarshal(body, &resData)
	if err != nil {
		return nil, fmt.Errorf("error parsing response data: %w", err)
	}

	page.Next = resData.Next
	page.Previous = resData.Previous

	areas := make([]string, 0, pageSize)
	for _, v := range resData.Results {
		areas = append(areas, v.Name)
	}

	return areas, nil
}

func QueryLocationAreaPokemon(locationArea string) ([]string, error) {
	url := apiLocationArea + locationArea

	body, err := apiGet(url)
	if err != nil {
		return nil, err
	}

	resData := resLocationArea{}
	err = json.Unmarshal(body, &resData)
	if err != nil {
		return nil, fmt.Errorf("error parsing response data: %w", err)
	}

	pokemon := make([]string, len(resData.PokemonEncounters))
	for i, v := range resData.PokemonEncounters {
		pokemon[i] = v.Pokemon.Name
	}

	return pokemon, nil
}
