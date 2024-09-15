package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	chirpydb "github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Make sure the database file exists
	db, err := chirpydb.NewDB("./database.json", *dbg)
	if err != nil {
		log.Fatalf("could not create database connection: %v", err)
	}

	state := newServerState(http.NewServeMux(), &apiConfig{}, db)

	serve := http.Server{
		Handler: state.Mux,
		Addr:    ":8080",
	}

	// Serve files in local directory without the /app prefix (namespace)
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	// Handle the entire /app/ path tree
	// This means not only /app, but also all subtrees under that path
	state.Mux.Handle("/app/", state.ApiCfg.middlewareMetricsInc(fileServerHandler))

	// Route the API using the multiplexer
	state.handleApi()

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
