package main

import (
	"fmt"
	"net/http"
)

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

func main() {
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	serve := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	// Serve files in local directory without the /app prefix (namespace)
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	// Handle the entire /app/ path tree
	// This means not only /app, but also all subtrees under that path
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

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

	// Channel for knowing when the server wants to stop
	stopChan := make(chan struct{})
	go func() {
		err := serve.ListenAndServe()
		if err != nil {
			// ListenAndServe never returns nil
			fmt.Println(err)
		}

		stopChan <- struct{}{}
	}()

	fmt.Println("Server up and running!")

	// Wait until server shuts down
	<-stopChan
}
