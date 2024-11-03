package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type UserIdResponse struct {
	UserId uuid.UUID `json:"user_id"`
}

type WebHookResponse struct {
	Event string         `json:"event"`
	Data  UserIdResponse `json:"data"`
}

func (cfg *apiConfig) upgradeUserToRed(w http.ResponseWriter, r *http.Request) {

	webhookParams := WebHookResponse{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&webhookParams); err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't decode parameters"}, http.StatusBadRequest)
		return
	}

	if webhookParams.Event != "user.upgraded" {
		respondWithJson(w, nil, http.StatusNoContent)
		return
	}

	_, err := cfg.db.UpgradeUserToRed(r.Context(), webhookParams.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJson(w, errorResponse{Error: "couldn't find user"}, http.StatusNotFound)
			return
		}
		respondWithJson(w, errorResponse{Error: "couldn't upgrade user"}, http.StatusInternalServerError)
		return
	}

	respondWithJson(w, nil, http.StatusNoContent)
}
