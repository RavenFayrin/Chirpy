package main

import (
	"Chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, convertUser(user))

}

func convertUser(dbUser database.User) (regUser User) {
	regUser.ID = dbUser.ID
	regUser.CreatedAt = dbUser.CreatedAt
	regUser.UpdatedAt = dbUser.UpdatedAt
	regUser.Email = dbUser.Email
	return regUser
}
