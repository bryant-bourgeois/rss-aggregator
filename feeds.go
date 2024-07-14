package main

import (
	"encoding/json"
	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) NewFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		Name string
		Url  string
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.Decode(&params)
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    u.ID,
	}
	feed, err := cfg.DB.CreateFeed(r.Context(), feedParams)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 201, feed)
}
