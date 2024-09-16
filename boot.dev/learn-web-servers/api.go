package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mosamadeeb/chirpy/internal/chirpydb"
)

type serverState struct {
	Mux    *http.ServeMux
	ApiCfg *apiConfig
	DB     *chirpydb.DB
}

func newServerState(mux *http.ServeMux, apiCfg *apiConfig, db *chirpydb.DB) serverState {
	return serverState{mux, apiCfg, db}
}

func (s serverState) handleApi() {
	// Admin namespace
	s.Mux.Handle("GET /admin/metrics", s.ApiCfg.adminMetricsHandler())

	// Metrics
	s.Mux.Handle("/api/reset", s.ApiCfg.resetHandler())

	// Readiness endpoint
	s.Mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Write status code before writing body
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	s.handleAuthApi()

	// CRUD endpoints
	s.handleChirpsApi()
	s.handleUsersApi()
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
		log.Printf("Error encoding %T response: %v\n", payload, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
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
