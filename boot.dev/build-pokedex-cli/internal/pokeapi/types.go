package pokeapi

type resLocationArea struct {
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
