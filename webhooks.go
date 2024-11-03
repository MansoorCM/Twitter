package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MansoorCM/Twitter/internal/auth"
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

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithJson(w, errorResponse{Error: "couldn't find api key"}, http.StatusUnauthorized)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithJson(w, errorResponse{Error: "invalid api key"}, http.StatusUnauthorized)
		return
	}

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

	_, err = cfg.db.UpgradeUserToRed(r.Context(), webhookParams.Data.UserId)
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
