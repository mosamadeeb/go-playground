package pokeapi

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
