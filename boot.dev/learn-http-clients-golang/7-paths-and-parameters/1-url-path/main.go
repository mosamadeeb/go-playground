package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	locations := getResources("/v1/courses_rest_api/learn-http/locations")
	fmt.Println("Locations:")
	logResources(locations)
	fmt.Println(" --- ")

	items := getResources("/v1/courses_rest_api/learn-http/items")
	fmt.Println("Items:")
	logResources(items)
	fmt.Println(" --- ")

	users := getResources("/v1/courses_rest_api/learn-http/users")
	fmt.Println("Users:")
	logResources(users)
}

func logResources(resources []map[string]any) {
	for _, resource := range resources {
		jsonResource, err := json.Marshal(resource)
		if err != nil {
			fmt.Println("Error marshalling resource:", err)
			continue
		}
		fmt.Printf(" - %s\n", jsonResource)
	}
}
