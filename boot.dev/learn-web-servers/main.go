package main

import (
	"fmt"
	"log"
	"net/http"

	chirpydb "github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func main() {
	apiCfg := apiConfig{}

	// Make sure the database file exists
	db, err := chirpydb.NewDB("./database.json")
	if err != nil {
		log.Fatalf("could not create database connection: %v", err)
	}

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

	handleApi(mux, &apiCfg, db)

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
