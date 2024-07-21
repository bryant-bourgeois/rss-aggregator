package main

import (
	"net/http"
	"strconv"

	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
)

func (cfg *apiConfig) GetPostsByUser(w http.ResponseWriter, r *http.Request, u database.User) {
	startRequested := r.URL.Query().Get("start")
	limitRequested := r.URL.Query().Get("limit")
	start, err := strconv.Atoi(startRequested)
	if err != nil {
		start = 0
	}
	limit, err := strconv.Atoi(limitRequested)
	if err != nil {
		limit = 100
	}
	query, err := cfg.DB.GetPostsByUserId(r.Context(), database.GetPostsByUserIdParams{
		ID:     u.ID,
		Offset: int32(start),
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, 500, "There was an error getting posts for your user...")
		return
	}
	posts := databasePostsToPosts(query)
	respondWithJSON(w, 200, posts)
	return
}
