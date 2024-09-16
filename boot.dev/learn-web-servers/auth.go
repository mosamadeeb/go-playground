package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mosamadeeb/chirpy/internal/chirpydb"
	"golang.org/x/crypto/bcrypt"
)

func (s serverState) handleAuthApi() {
	s.Mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		var loginReq struct {
			Email            string `json:"email"`
			Password         string `json:"password"`
			ExpiresInSeconds int    `json:"expires_in_seconds"`
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

		timeoutSec := loginReq.ExpiresInSeconds
		timeoutDefault := 24 * int(time.Hour/time.Second)
		if timeoutSec <= 0 || timeoutSec > timeoutDefault {
			timeoutSec = timeoutDefault
		}

		jwtToken, err := createJWT(user.Id, timeoutSec, s.ApiCfg.jwtSecret)
		if err != nil {
			log.Printf("Error creating JWT: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
			Token string `json:"token"`
		}{user.Id, user.Email, jwtToken})
	})
}

func createJWT(userId, expiresInSeconds int, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(expiresInSeconds))),
		Subject:   strconv.Itoa(userId),
	})

	return token.SignedString([]byte(jwtSecret))
}
