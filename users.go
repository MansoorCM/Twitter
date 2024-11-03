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
	Id           string `json:"id"`
	Created_at   string `json:"created_at"`
	Updated_at   string `json:"updated_at"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't find JWT"}, http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't validate JWT"}, http.StatusUnauthorized)
		return
	}

	user := User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't decode parameters"}, http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "invalid password"}, http.StatusBadRequest)
		return
	}

	userDb, err := cfg.db.UpdateUser(r.Context(),
		database.UpdateUserParams{
			Email:          user.Email,
			HashedPassword: hashedPassword,
			ID:             userId})

	if err != nil {
		respondWithJson(w, errorResponse{Error: "failed to update user"}, http.StatusInternalServerError)
		return
	}

	userResponse := UserResponse{Id: userDb.ID.String(),
		Created_at: userDb.CreatedAt.String(),
		Updated_at: userDb.UpdatedAt.String(),
		Email:      userDb.Email}

	respondWithJson(w, userResponse, http.StatusOK)
}
