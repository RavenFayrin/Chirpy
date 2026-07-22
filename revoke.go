package main

import (
	"errors"
	"net/http"
	"strings"
)


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusInternalServerError, "malformed authorization header", errors.New("no auth header included in request"))
		return
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		respondWithError(w, http.StatusInternalServerError, "malformed authorization header", errors.New("no auth header included in request"))
		return
	}

	err := cfg.db.RevokeRefreshToken(r.Context(), splitAuth[1])
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token does not match or exist", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}