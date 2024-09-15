package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func (s serverState) handleChirpsApi() {
	s.Mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
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

		chirp, err := s.DB.CreateChirp(cleanChirp(chirpReq.Body))
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

		respondWithJSON(w, http.StatusOK, chirps)
	})

	s.Mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chirps, err := s.DB.GetChirps()
		if err != nil {
			log.Printf("Error loading chirps from database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// The chirps are sorted by ID so we can do a binary search
		index, ok := slices.BinarySearchFunc(chirps, chirpydb.Chirp{Id: chirpId}, func(a, b chirpydb.Chirp) int { return a.Id - b.Id })

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		respondWithJSON(w, http.StatusOK, chirps[index])
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
