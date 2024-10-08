package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
	"golang.org/x/crypto/bcrypt"
)

// Mirrors User but with password removed
type userRes struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func createUserRes(u chirpydb.User) userRes {
	return userRes{
		u.Id,
		u.Email,
		u.IsChirpyRed,
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
			if errors.Is(err, chirpydb.ErrExists) {
				// Email already used
				w.WriteHeader(http.StatusConflict)
				return
			}

			log.Printf("Error saving user to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, createUserRes(user))
	})

	s.Mux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, r *http.Request) {
		token, err := authenticateJWT(r, s.ApiCfg.jwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		idStr, err := token.Claims.GetSubject()
		if err != nil {
			log.Printf("Could not get user ID from JWT: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Could not get user ID from JWT: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		user, err := s.DB.UpdateUser(userId, userReq.Email, string(encPassword))
		if err != nil {
			if errors.Is(err, chirpydb.ErrNotExist) {
				// Ah yes, user must have deleted their account and *then* proceeded to update their credentials
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if errors.Is(err, chirpydb.ErrExists) {
				// Email already used
				w.WriteHeader(http.StatusConflict)
				return
			}

			log.Printf("Error updating user in database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, createUserRes(user))
	})
}
