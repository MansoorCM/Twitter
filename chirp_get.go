package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirpsFromDb, err := cfg.db.GetAllChirps(r.Context())
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

	respondWithJson(w, chirps, http.StatusOK)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

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
