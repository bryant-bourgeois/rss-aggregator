package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type messageResponse struct {
	Message string `json:"message"`
}

func (cfg *apiConfig) NewUser(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Name string
	}
	reqData := req{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&reqData)
	if err != nil {
		respondWithJSON(w, 400, messageResponse{Message: "Invalid request"})
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqData.Name,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (cfg *apiConfig) GetUser(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" || apiKey[:6] != "ApiKey" {
		respondWithJSON(w, 401, messageResponse{Message: "Need to send a 'Authorization: ApiKey API_KEY' header."})
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey[7:])
	if err != nil {
		respondWithError(w, 401, "User not found or api key invalid")
		return
	}
	respondWithJSON(w, 200, databaseUserToUser(user))
}
