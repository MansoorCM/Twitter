package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Chirp struct {
	Body string `json:"body"`
}

type CleanChirp struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
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

	chirpAfterCleaning := replaceProfaneWords(chirp.Body)
	resp := CleanChirp{CleanedBody: chirpAfterCleaning}
	respondWithJson(w, resp, http.StatusOK)
}

func replaceProfaneWords(s string) string {
	profaneWords := getProfaneWords()

	words := strings.Split(s, " ")

	for i, word := range words {
		lowercaseWord := strings.ToLower(word)
		if _, ok := profaneWords[lowercaseWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func getProfaneWords() map[string]struct{} {
	profaneWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	return profaneWords
}
