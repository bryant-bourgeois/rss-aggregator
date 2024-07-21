package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bryant-bourgeois/rss-aggregator/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) NewFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
		return
	}

	// Some valid RSS feed urls don't serve an xml document directly
	// if params.Url[len(params.Url)-4:] != ".xml" {
	// 	respondWithError(w, 400, "Supplied url was not to an RSS feed")
	// 	return
	// }

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
	follow, err := cfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		UserID:    u.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	type newFeedAndSubscription struct {
		Feed   Feed       `json:"feed"`
		Follow FeedFollow `json:"feed_follow"`
	}
	respondWithJSON(w, 201, newFeedAndSubscription{
		Feed:   databaseFeedToFeed(feed),
		Follow: databaseFeedFollowToFeedFollow(follow),
	})
}

func (cfg *apiConfig) GetFeeds(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := cfg.DB.ListFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	userFeeds := make([]Feed, 0)
	for _, val := range dbFeeds {
		userFeeds = append(userFeeds, databaseFeedToFeed(val))
	}
	respondWithJSON(w, 200, userFeeds)
}
