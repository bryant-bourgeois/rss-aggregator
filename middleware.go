package main

import (
	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
	"net/http"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		handler(w, r, user)
	})
}
