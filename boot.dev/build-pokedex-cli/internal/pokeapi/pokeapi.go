package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	apiBase         = "https://pokeapi.co/api/v2/"
	apiLocationArea = apiBase + "location-area/"
	pageSize        = 20
)

type PageConfig struct {
	Next     *string
	Previous *string
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

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("api get request error: %w", err)
	}

	// Important: defer closing the body reader *after* checking for the error
	// If the error is nil, the body will be already closed
	// If the error is not nil, closing the body will allow us to reuse the connection
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	resData := resLocationArea{}
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
