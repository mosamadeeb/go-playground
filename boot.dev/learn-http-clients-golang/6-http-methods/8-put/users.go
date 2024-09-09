package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func updateUser(baseURL, id, apiKey string, data User) (User, error) {
	fullURL := baseURL + "/" + id

	marshalled, err := json.Marshal(data)
	if err != nil {
		return User{}, err
	}

	req, err := http.NewRequest("PUT", fullURL, bytes.NewBuffer(marshalled))
	if err != nil {
		return User{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()

	var updatedUser User
	if err = json.NewDecoder(res.Body).Decode(&updatedUser); err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func getUserById(baseURL, id, apiKey string) (User, error) {
	fullURL := baseURL + "/" + id

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Set("X-API-Key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()

	var user User
	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

// don't touch below this line

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
