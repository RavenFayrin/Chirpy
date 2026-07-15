package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	split_msg := strings.Split(params.Body, " ")
	for i, msg := range split_msg {
		if (strings.ToLower(msg) == "kerfuffle") || (strings.ToLower(msg) == "sharbert") || (strings.ToLower(msg) == "fornax") {
			split_msg[i] = "****"
		}
	}
	params.Body = strings.Join(split_msg, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_Body:	params.Body,
	})
}
