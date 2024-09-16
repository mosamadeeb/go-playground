package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
	"golang.org/x/crypto/bcrypt"
)

// Mirrors User but with password removed
type userRes struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func createUserRes(u chirpydb.User) userRes {
	return userRes{
		u.Id,
		u.Email,
	}
}

func (s serverState) handleUsersApi() {
	s.Mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		var userReq struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
			log.Printf("Error decoding user body: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		encPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing user password: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.DB.CreateUser(userReq.Email, string(encPassword))
		if err != nil {
			log.Printf("Error saving user to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, createUserRes(user))
	})
}
