package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func (s serverState) handleChirpsApi() {
	s.Mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
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

		var chirpReq chirpydb.Chirp
		if err := json.NewDecoder(r.Body).Decode(&chirpReq); err != nil {
			log.Printf("Error decoding chirp body: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(chirpReq.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		chirp, err := s.DB.CreateChirp(cleanChirp(chirpReq.Body), userId)
		if err != nil {
			log.Printf("Error saving chirp to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, chirp)
	})

	s.Mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		chirps, err := s.DB.GetChirps()
		if err != nil {
			log.Printf("Error loading chirps from database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if authorIdStr := r.URL.Query().Get("author_id"); authorIdStr != "" {
			authorId, err := strconv.Atoi(authorIdStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			filteredChirps := make([]chirpydb.Chirp, 0, len(chirps))
			for _, c := range chirps {
				if c.AuthorId == authorId {
					filteredChirps = append(filteredChirps, c)
				}
			}

			chirps = filteredChirps
		}

		respondWithJSON(w, http.StatusOK, chirps)
	})

	s.Mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chirp, err := s.DB.GetChirp(chirpId)
		if err != nil {
			if errors.Is(err, chirpydb.ErrNotExist) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			log.Printf("Error loading chirp from database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, chirp)
	})

	s.Mux.HandleFunc("DELETE /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
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

		chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chirp, err := s.DB.GetChirp(chirpId)
		if err != nil {
			if errors.Is(err, chirpydb.ErrNotExist) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			log.Printf("Error loading chirp from database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userId != chirp.AuthorId {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err = s.DB.DeleteChirp(chirpId)
		if err != nil && !errors.Is(err, chirpydb.ErrNotExist) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func cleanChirp(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(body, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
