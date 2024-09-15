package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func handleApi(mux *http.ServeMux, apiCfg *apiConfig, db *chirpydb.DB) {
	// Admin namespace
	mux.Handle("GET /admin/metrics", apiCfg.adminMetricsHandler())

	// Metrics
	mux.Handle("/api/reset", apiCfg.resetHandler())

	// Readiness endpoint
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Write status code before writing body
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chirp CRUD endpoints
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		var chirpReq chirpydb.Chirp
		if err := json.NewDecoder(r.Body).Decode(&chirpReq); err != nil {
			log.Printf("Error decoding chirp body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(chirpReq.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		chirp, err := db.CreateChirp(cleanChirp(chirpReq.Body))
		if err != nil {
			log.Printf("Error saving chirp to database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, chirp)
	})

	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetChirps()
		if err != nil {
			log.Printf("Error loading chirps from database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
	})

	mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chirps, err := db.GetChirps()
		if err != nil {
			log.Printf("Error loading chirps from database: %v", err)
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

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		var userReq chirpydb.User
		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
			log.Printf("Error decoding chirp body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := db.CreateUser(userReq.Email)
		if err != nil {
			log.Printf("Error saving user to database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResp struct {
		ErrorMsg string `json:"error"`
	}

	if len(msg) == 0 {
		w.WriteHeader(code)
	} else {
		respondWithJSON(w, code, errorResp{msg})
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding %T response: %s", payload, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
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

type apiConfig struct {
	fileserverHits int
}

// A simple middleware that inserts a handler in between
// This allows us to do something before the next handler is used
func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

// We can return a handler by wrapping a function with an appropriate signature with http.HandlerFunc()
func (c *apiConfig) resetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits = 0
	})
}

func (c *apiConfig) adminMetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, c.fileserverHits)))
	})
}
