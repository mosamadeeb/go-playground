package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	chirpydb "github.com/mosamadeeb/chirpy/internal/chirpydb"
)

func main() {
	godotenv.Load()

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Make sure the database file exists
	db, err := chirpydb.NewDB("./database.json", *dbg)
	if err != nil {
		log.Fatalf("could not create database connection: %v\n", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	state := newServerState(http.NewServeMux(), &apiConfig{jwtSecret: jwtSecret}, db)

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
