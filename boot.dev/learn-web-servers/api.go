package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handleApi(mux *http.ServeMux, apiCfg *apiConfig) {
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

	// Chirp validation endpoint
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		var chirpBody struct {
			Body string `json:"body"`
		}

		if err := json.NewDecoder(r.Body).Decode(&chirpBody); err != nil {
			log.Printf("Error decoding chirp body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type errorResp struct {
			ErrorMsg string `json:"error"`
		}

		if len(chirpBody.Body) > 140 {
			resp, err := json.Marshal(errorResp{"Chirp is too long"})
			if err != nil {
				log.Printf("Error encoding error response: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}

		type cleanedResp struct {
			Body string `json:"cleaned_body"`
		}

		badWords := []string{"kerfuffle", "sharbert", "fornax"}
		words := strings.Split(chirpBody.Body, " ")
		for i, word := range words {
			if slices.Contains(badWords, strings.ToLower(word)) {
				words[i] = "****"
			}
		}

		resp, err := json.Marshal(cleanedResp{strings.Join(words, " ")})
		if err != nil {
			log.Printf("Error encoding valid response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	})
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
