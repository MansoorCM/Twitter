package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't decode parameters"}, http.StatusInternalServerError)
		return
	}

	dbResponse, err := cfg.db.CreateUser(r.Context(), user.Email)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "failed to create user"}, http.StatusInternalServerError)
		return
	}

	respondWithJson(w, dbResponse, http.StatusCreated)
}
