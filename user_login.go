package main

import (
	"encoding/json"
	"net/http"

	"github.com/MansoorCM/Twitter/internal/auth"
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

	userResponse := UserResponse{Id: user.ID.String(),
		Created_at: user.CreatedAt.String(),
		Updated_at: user.UpdatedAt.String(),
		Email:      user.Email}

	respondWithJson(w, userResponse, http.StatusOK)
}
