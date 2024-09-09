package main

import (
	"fmt"
	"log"
)

type User struct {
	CharacterName string `json:"characterName"`
	Class         string `json:"class"`
	Level         int    `json:"level"`
	PvpEnabled    bool   `json:"pvpEnabled"`
	User          struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Age      int    `json:"age"`
	} `json:"user"`
}

func logUsers(users []User) {
	for _, user := range users {
		fmt.Printf("Character name: %s, Class: %s, Level: %d, User: %s\n", user.CharacterName, user.Class, user.Level, user.User.Name)
	}
}

func main() {
	baseURL := "https://api.boot.dev/v1/courses_rest_api/learn-http/users"

	users, err := getUsers(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	logUsers(users)

}
