package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/mosamadeeb/chirpy/internal/database"
)

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("could not create database connection: %v\n", err)
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{db: dbQueries}

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

	handleApi(mux, &apiCfg)

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
