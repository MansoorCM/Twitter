package main

import (
	"net/http"

	"github.com/MansoorCM/Twitter/internal/auth"
	"github.com/MansoorCM/Twitter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	authorId := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	var chirpsFromDb []database.Chirp
	var err error

	authorUUID, err := uuid.Parse(authorId)

	if err != nil || authorId == "" {
		chirpsFromDb, err = cfg.db.GetAllChirps(r.Context())
	} else {
		chirpsFromDb, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
	}

	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't retrieve chirps"}, http.StatusInternalServerError)
		return
	}

	chirps := make([]ChirpResponse, len(chirpsFromDb))
	for i, dbChirp := range chirpsFromDb {
		chirp := ChirpResponse{ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body}
		chirps[i] = chirp
	}

	if sortOrder == "desc" {
		reverseList(chirps)
	}

	respondWithJson(w, chirps, http.StatusOK)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")

	chirp_uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "Invalid uuid"}, http.StatusBadRequest)
		return
	}

	chirpFromDb, err := cfg.db.GetChirp(r.Context(), chirp_uuid)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't retrieve chirp from db"}, http.StatusNotFound)
		return
	}

	chirp := ChirpResponse{ID: chirpFromDb.ID,
		CreatedAt: chirpFromDb.CreatedAt,
		UpdatedAt: chirpFromDb.UpdatedAt,
		UserID:    chirpFromDb.UserID,
		Body:      chirpFromDb.Body}

	respondWithJson(w, chirp, http.StatusOK)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't find JWT"}, http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't validate JWT"}, http.StatusForbidden)
		return
	}

	id := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(id)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "invalid uuid"}, http.StatusBadRequest)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't find chirp in DB"}, http.StatusNotFound)
		return
	}
	if dbChirp.UserID != userId {
		respondWithJson(w, errorResponse{Error: "forbidden operation"}, http.StatusForbidden)
		return
	}

	_, err = cfg.db.DeleteChirp(r.Context(),
		database.DeleteChirpParams{
			ID: chirpUUID, UserID: userId,
		})
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't delete chirp from DB"}, http.StatusInternalServerError)
		return
	}

	respondWithJson(w, nil, http.StatusNoContent)
}

func reverseList[T any](items []T) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}
