package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, resp interface{}, status int) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResp)
}
