package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	serve := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	// Channel for knowing when the server wants to stop
	stopChan := make(chan struct{})
	go func() {
		err := serve.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

		stopChan <- struct{}{}
	}()

	fmt.Println("Server up and running!")

	// Wait until server shuts down
	<-stopChan
}
