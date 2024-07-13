package main

import (
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, 200, resp{Status: "ok"})
}

func healthCheckError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}
