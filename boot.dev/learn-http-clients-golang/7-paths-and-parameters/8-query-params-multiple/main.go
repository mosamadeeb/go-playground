package main

import "strconv"

func lootTreasure(baseURL, chestRarity string) []Item {
	limit := 0
	switch chestRarity {
	case "Common":
		limit = 1
	case "Rare":
		limit = 3
	case "Legendary":
		limit = 5
	}

	fullURL := baseURL + "?sort=quality&limit=" + strconv.Itoa(limit)
	return getItems(fullURL)
}
