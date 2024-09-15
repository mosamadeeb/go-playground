package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func (s serverState) handleUsersApi() {
	s.Mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		var userReq chirpydb.User
		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
			log.Printf("Error decoding chirp body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.DB.CreateUser(userReq.Email)
		if err != nil {
			log.Printf("Error saving user to database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	})
}
