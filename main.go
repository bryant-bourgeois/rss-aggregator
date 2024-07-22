package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Couldn't load environment variable file: %s\n", err)
		os.Exit(1)
	}
	refreshInterval, err := strconv.Atoi(os.Getenv("DB_FETCH_INTERVAL_SECONDS"))
	if err != nil {
		refreshInterval = 60
	}
	refreshAmount, err := strconv.Atoi(os.Getenv("DB_FETCH_COUNT"))
	if err != nil {
		refreshAmount = 3
	}
	cfg := apiConfig{
		FeedRefreshAmount:         refreshAmount,
		FeedUpdateIntervalSeconds: refreshInterval,
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
	mux.Handle("GET /v1/users", cfg.middlewareAuth(cfg.GetUser))

	mux.HandleFunc("GET /v1/feeds", cfg.GetFeeds)
	mux.Handle("POST /v1/feeds", cfg.middlewareAuth(cfg.NewFeed))

	mux.Handle("GET /v1/feed_follows", cfg.middlewareAuth(cfg.GetUserFeedFollows))
	mux.Handle("POST /v1/feed_follows", cfg.middlewareAuth(cfg.FollowFeed))
	mux.Handle("DELETE /v1/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.DeleteFeedFollow))

	mux.Handle("GET /v1/posts", cfg.middlewareAuth(cfg.GetPostsByUser))

	go RefreshFeeds(cfg)
	fmt.Printf("Starting server on %s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
