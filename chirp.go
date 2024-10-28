package main

import (
	"encoding/json"
	"net/http"
)

type Chirp struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type returnVals struct {
	Valid bool `json:"valid"`
}

func validateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}

	if err := decoder.Decode(&chirp); err != nil {
		errResp := errorResponse{Error: "Something went wrong"}
		respondWithJson(w, errResp, http.StatusInternalServerError)
		return
	}

	if len(chirp.Body) > 140 {
		errResp := errorResponse{Error: "Chirp is too long"}
		respondWithJson(w, errResp, http.StatusBadRequest)
		return
	}

	resp := returnVals{Valid: true}
	respondWithJson(w, resp, http.StatusOK)
}
