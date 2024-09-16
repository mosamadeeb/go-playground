package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
	"golang.org/x/crypto/bcrypt"
)

func (s serverState) handleAuthApi() {
	s.Mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		var loginReq struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			log.Printf("Error decoding user body: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loginErrMsg := "Incorrect email or password"

		user, err := s.DB.GetUserByEmail(loginReq.Email)
		if err != nil {
			if errors.Is(err, chirpydb.ErrNotExist) {
				// Returning 401 if user is not found
				respondWithError(w, http.StatusUnauthorized, loginErrMsg)
				return
			} else {
				log.Printf("Error fetching user from database: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
			respondWithError(w, http.StatusUnauthorized, loginErrMsg)
			return
		}

		respondWithJSON(w, http.StatusOK, createUserRes(user))
	})
}
