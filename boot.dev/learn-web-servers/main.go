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

	// Serve files in local directory
	mux.Handle("/", http.FileServer(http.Dir(".")))

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
