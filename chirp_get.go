package main

import "net/http"

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
