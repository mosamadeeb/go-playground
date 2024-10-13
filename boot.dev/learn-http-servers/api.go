package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mosamadeeb/chirpy/internal/database"
)

func handleApi(mux *http.ServeMux, apiCfg *apiConfig) {
	// Admin namespace
	mux.Handle("GET /admin/metrics", apiCfg.adminMetricsHandler())

	// Metrics
	mux.Handle("/admin/reset", apiCfg.resetHandler())

	// Readiness endpoint
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Write status code before writing body
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chirp endpoints
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		var chirpBody struct {
			Body   string `json:"body"`
			UserId string `json:"user_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&chirpBody); err != nil {
			log.Printf("Error decoding chirp body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(chirpBody.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		cleanedBody := cleanChirp(chirpBody.Body)
		userId, err := uuid.Parse(chirpBody.UserId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chirp, err := apiCfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanedBody, UserID: userId})
		if err != nil {
			log.Printf("Error creating chirp in database: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, Chirp{
			chirp.ID,
			chirp.CreatedAt,
			chirp.UpdatedAt,
			chirp.Body,
			chirp.UserID.String(),
		})
	})

	// User endpoints
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Printf("Error decoding request body: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := apiCfg.db.CreateUser(r.Context(), body.Email)
		if err != nil {
			log.Printf("Error creating user in database: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email})
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
	db             *database.Queries
	platform       string
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
		if c.platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		c.fileserverHits = 0

		// Delete everything from the DB!!
		c.db.DeleteAllUsers(r.Context())

		w.WriteHeader(http.StatusOK)
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

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
}
