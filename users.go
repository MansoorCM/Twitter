package main

import (
	"encoding/json"
	"net/http"

	"github.com/MansoorCM/Twitter/internal/auth"
	"github.com/MansoorCM/Twitter/internal/database"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't decode parameters"}, http.StatusInternalServerError)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "Invalid password"}, http.StatusInternalServerError)
		return
	}

	dbResponse, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithJson(w, errorResponse{Error: "failed to create user"}, http.StatusInternalServerError)
		return
	}

	userResponse := UserResponse{Id: dbResponse.ID.String(),
		Created_at: dbResponse.CreatedAt.String(),
		Updated_at: dbResponse.UpdatedAt.String(),
		Email:      dbResponse.Email}

	respondWithJson(w, userResponse, http.StatusCreated)
}
