package main

import (
	"Chirpy/internal/auth"
	"errors"
	"net/http"
	"strings"
	"time"
)


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token        string `json:"token"`
	}

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

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), splitAuth[1])
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "refresh token does not match or exist", errors.New("no auth header included in request"))
		return
	}

	expirationTime := time.Hour

	accessToken, err := auth.MakeJWT(
		refreshToken.UserID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
