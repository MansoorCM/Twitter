package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/MansoorCM/Twitter/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

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

	chirp.Body = replaceProfaneWords(chirp.Body)

	dbResponse, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirp.Body,
		UserID: chirp.UserID,
	})

	if err != nil {
		respondWithJson(w, errorResponse{Error: "failed to create chirp"}, http.StatusInternalServerError)
		return
	}

	chirpResponse := ChirpResponse{ID: dbResponse.ID,
		CreatedAt: dbResponse.CreatedAt,
		UpdatedAt: dbResponse.UpdatedAt,
		UserID:    dbResponse.UserID,
		Body:      dbResponse.Body}

	respondWithJson(w, chirpResponse, http.StatusCreated)
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
