package main

import (
	"fmt"
)

func main() {
	url := "https://api.boot.dev/v1/courses_rest_api/learn-http/locations/52fdfc07-2182-454f-963f-5f0f9a621d72"
	apiKey := generateKey()

	oldLocation, err := getLocationResponse(apiKey, url)
	if err != nil {
		fmt.Println("Error getting old location:", err)
		return
	}
	fmt.Println("Got old location:")
	fmt.Printf("- name: %s\n", oldLocation.Name)
	fmt.Printf("- recommendedLevel: %d\n", oldLocation.RecommendedLevel)
	fmt.Println("--------------------------------")

	newLocationData := Location{
		Discovered:       false,
		ID:               "52fdfc07-2182-454f-963f-5f0f9a621d72",
		Name:             "Bloodstone Swamp",
		RecommendedLevel: 10,
	}

	if err := putLocation(apiKey, url, newLocationData); err != nil {
		fmt.Println("Error updating location:", err)
		return
	}
	fmt.Println("Location updated!")
	fmt.Println("---")

	newLocation, err := getLocationResponse(apiKey, url)
	if err != nil {
		fmt.Println("Error getting new location:", err)
		return
	}
	fmt.Println("Got new location:")
	fmt.Printf("- name: %s\n", newLocation.Name)
	fmt.Printf("- recommendedLevel: %d\n", newLocation.RecommendedLevel)
	fmt.Println("--------------------------------")
}
