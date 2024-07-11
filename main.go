package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Couldn't load environment variable file: %s\n", err)
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	mux.HandleFunc("GET /v1/healthz", healthCheck)
	mux.HandleFunc("GET /v1/err", healthCheckError)

	fmt.Printf("Starting server on %s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
