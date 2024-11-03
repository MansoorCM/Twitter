package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MansoorCM/Twitter/internal/auth"
	"github.com/MansoorCM/Twitter/internal/database"
)

type LoginDetails struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) userLogin(w http.ResponseWriter, r *http.Request) {
	loginDetails := LoginDetails{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginDetails); err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't decode parameters"}, http.StatusInternalServerError)
		return
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

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't create JWT"}, http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		errorMsg := fmt.Sprintf("Couldn't create refresh token %v", err)
		respondWithJson(w, errorResponse{Error: errorMsg}, http.StatusInternalServerError)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		errorMsg := fmt.Sprintf("Couldn't save refresh token %v", err)
		respondWithJson(w, errorResponse{Error: errorMsg}, http.StatusInternalServerError)
		return
	}

	userResponse := UserResponse{Id: user.ID.String(),
		Created_at:   user.CreatedAt.String(),
		Updated_at:   user.UpdatedAt.String(),
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken}

	respondWithJson(w, userResponse, http.StatusOK)
}
