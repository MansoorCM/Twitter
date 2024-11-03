package main

import (
	"net/http"
	"time"

	"github.com/MansoorCM/Twitter/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type refreshResponse struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't find token"}, http.StatusBadRequest)
		return
	}

	refreshTokenDb, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't get user for refresh token"}, http.StatusUnauthorized)
		return
	}

	newToken, err := auth.MakeJWT(refreshTokenDb.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't create JWT"}, http.StatusUnauthorized)
		return
	}

	respondWithJson(w, refreshResponse{Token: newToken}, http.StatusOK)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't find token"}, http.StatusBadRequest)
		return
	}

	_, err = cfg.db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't revoke session"}, http.StatusUnauthorized)
		return
	}

	respondWithJson(w, nil, http.StatusNoContent)
}
