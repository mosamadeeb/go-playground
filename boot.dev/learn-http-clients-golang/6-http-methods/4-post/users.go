package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func createUser(url, apiKey string, data User) (User, error) {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return User{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalled))
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

	var newUser User
	if err = json.NewDecoder(res.Body).Decode(&newUser); err != nil {
		return User{}, err
	}

	return newUser, nil
}

// Don't touch below this line

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

func getUsers(url, apiKey string) ([]User, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []User{}, err
	}

	req.Header.Set("X-API-Key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []User{}, err
	}
	defer res.Body.Close()

	var users []User
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&users)
	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func logUsers(users []User) {
	for _, user := range users {
		fmt.Printf("Character name: %s, Class: %s, Level: %d, User: %s\n", user.CharacterName, user.Class, user.Level, user.User.Name)
	}
}
