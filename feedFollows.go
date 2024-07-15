package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"sort"
	"github.com/bryant-bourgeois/rss-aggregator/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) FollowFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
		return
	}
	dbFollow, err := cfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID: uuid.New(),
		UserID: u.ID,
		FeedID: params.FeedId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	userFollow := databaseFeedFollowToFeedFollow(dbFollow)
	respondWithJSON(w, 201, userFollow)
}

func (cfg *apiConfig) DeleteFeedFollow(w http.ResponseWriter, r *http.Request, u database.User) {
	targetFeed := r.PathValue("feedFollowID")
	targetFeedUUID, err := uuid.Parse(targetFeed)
	if err != nil {
		respondWithError(w, 400, "invalid feed ID")
		return
	}
	dbFeedFollow, err := cfg.DB.GetFeedFollow(r.Context(), targetFeedUUID)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if dbFeedFollow.UserID != u.ID {
		respondWithError(w, 403, "Forbidden: you are not authorized to unfollow this feed")
		return
	}
	err = cfg.DB.UnfollowFeed(r.Context(), targetFeedUUID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, messageResponse{Message: fmt.Sprintf("Successfully unfollowed feed %v", targetFeedUUID)})	
}

func (cfg *apiConfig) GetUserFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	follows, err := cfg.DB.ListFeedFollows(r.Context(), u.ID)
	if err != nil  {
		respondWithError(w, 500, "Internal server error")
		return
	}
	userFollows := make([]FeedFollow, 0)
	for _, val := range follows {
		userFollows = append(userFollows, databaseFeedFollowToFeedFollow(val))
	}
	sort.Slice(userFollows, func(i, j int ) bool {
		return userFollows[i].CreatedAt.Before(userFollows[j].CreatedAt)
	})
	respondWithJSON(w, 200, userFollows)
	
}
