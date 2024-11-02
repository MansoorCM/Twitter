package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MansoorCM/Twitter/internal/auth"
)

type LoginDetails struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) userLogin(w http.ResponseWriter, r *http.Request) {
	loginDetails := LoginDetails{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginDetails); err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't decode parameters"}, http.StatusInternalServerError)
		return
	}

	expire_time := time.Hour
	if loginDetails.ExpiresInSeconds > 0 && loginDetails.ExpiresInSeconds < 3600 {
		expire_time = time.Second * time.Duration(loginDetails.ExpiresInSeconds)
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), loginDetails.Email)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "Incorrect Email or password"}, http.StatusUnauthorized)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, loginDetails.Password)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "Incorrect Email or password"}, http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expire_time)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't create JWT"}, http.StatusInternalServerError)
		return
	}

	userResponse := UserResponse{Id: user.ID.String(),
		Created_at: user.CreatedAt.String(),
		Updated_at: user.UpdatedAt.String(),
		Email:      user.Email,
		Token:      token}

	respondWithJson(w, userResponse, http.StatusOK)
}
