package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeMembership(w http.ResponseWriter, r *http.Request) {
	type webhookData struct {
    UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  webhookData `json:"data"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeMembershipById(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't upgrade membership", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
