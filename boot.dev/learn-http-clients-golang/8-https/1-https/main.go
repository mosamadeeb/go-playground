package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const URL = "https://api.boot.dev/v1/courses_rest_api/learn-http/users"

// don't touch below this line

func getUserById(baseURL, id string) (User, error) {
	fullURL := baseURL + "/" + id
	if err := errIfNotHTTPS(fullURL); err != nil {
		return User{}, err
	}

	res, err := http.Get(fullURL)
	if err != nil {
		return User{}, err
	}

	var user User
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func main() {
	uuid := "2f8282cb-e2f9-496f-b144-c0aa4ced56db"

	user, err := getUserById(URL, uuid)
	if err != nil {
		fmt.Println(err)
		return
	}

	logUser(user)
}

type User struct {
	CharacterName string `json:"characterName"`
	Class         string `json:"class"`
	ID            string `json:"id"`
	Level         int    `json:"level"`
	PvpEnabled    bool   `json:"pvpEnabled"`
	User          struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Age      int    `json:"age"`
	} `json:"user"`
}

func logUser(user User) {
	fmt.Printf("User uuid: %s, Character Name: %s, Class: %s, Level: %d, PVP Status: %t, User name: %s\n",
		user.ID, user.CharacterName, user.Class, user.Level, user.PvpEnabled, user.User.Name)
}

func errIfNotHTTPS(URL string) error {
	url, err := url.Parse(URL)
	if err != nil {
		return err
	}
	if url.Scheme != "https" {
		return fmt.Errorf("URL scheme is not HTTPS: %s", URL)
	}
	return nil
}
