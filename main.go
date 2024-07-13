package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	cfg := apiConfig{}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Couldn't load environment variable file: %s\n", err)
		os.Exit(1)
	}

	dbConn := os.Getenv("DB_CONN")
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		fmt.Printf("Couldn't connect to database: %s\n", err)
		os.Exit(1)
	}
	dbQueies := database.New(db)
	cfg.DB = dbQueies

	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.HandleFunc("GET /v1/healthz", healthCheck)
	mux.HandleFunc("GET /v1/err", healthCheckError)

	mux.HandleFunc("POST /v1/users", cfg.NewUser)
	mux.HandleFunc("GET /v1/users", cfg.GetUser)

	fmt.Printf("Starting server on %s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
