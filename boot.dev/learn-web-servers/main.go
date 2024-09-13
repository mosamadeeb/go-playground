package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	serve := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	// Serve files in local directory without the /app prefix
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// Readiness endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
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
