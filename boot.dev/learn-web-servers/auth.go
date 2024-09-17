package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mosamadeeb/chirpy/internal/chirpydb"
	"golang.org/x/crypto/bcrypt"
)

const jwtDefaultTimeout = time.Hour

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

		expiresIn := time.Duration(loginReq.ExpiresInSeconds) * time.Second
		if expiresIn <= 0 || expiresIn > jwtDefaultTimeout {
			expiresIn = jwtDefaultTimeout
		}

		jwtToken, err := createJWT(user.Id, expiresIn, s.ApiCfg.jwtSecret)
		if err != nil {
			log.Printf("Error creating JWT: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		refreshToken, err := s.DB.AddRefreshToken(user.Id, time.Now().AddDate(0, 0, 60).UTC())
		if err != nil {
			log.Printf("Error creating refresh token: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Id           int    `json:"id"`
			Email        string `json:"email"`
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}{user.Id, user.Email, jwtToken, refreshToken.Token})
	})

	s.Mux.HandleFunc("POST /api/refresh", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userId, err := s.DB.CheckRefreshToken(tokenString)
		if err != nil {
			if errors.Is(err, chirpydb.ErrNotExist) || errors.Is(err, chirpydb.ErrExpired) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Printf("Error checking refresh token in DB: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jwtToken, err := createJWT(userId, jwtDefaultTimeout, s.ApiCfg.jwtSecret)
		if err != nil {
			log.Printf("Error creating JWT: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Token string `json:"token"`
		}{jwtToken})
	})

	s.Mux.HandleFunc("POST /api/revoke", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := s.DB.RevokeRefreshToken(tokenString)
		if err != nil && !errors.Is(err, chirpydb.ErrNotExist) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func createJWT(userId int, expiresIn time.Duration, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   strconv.Itoa(userId),
	})

	return token.SignedString([]byte(jwtSecret))
}

// TODO: Instead of authenticating in a function that takes a request, we can use a middleware
func authenticateJWT(r *http.Request, jwtSecret string) (*jwt.Token, error) {
	tokenString, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		return nil, fmt.Errorf("unexpected authorization header format")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
