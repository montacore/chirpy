package main

import (
	"net/http"
	"log"
	"github.com/montacore/chirpy/internal/auth"
	"github.com/montacore/chirpy/internal/database"
	"encoding/json"
)


//TEST: Endpoint: "POST /api/chirps"
func (cfg *apiConfig)handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body 	string		`json:"body"`
	}
	// TODO: authorization first and foremost
	key, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	userID, err := auth.ValidateJWT(key, cfg.secretKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	// TODO: Decode r.Body into parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		RespondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
		
		return
	}

	// TODO: Check length and call respondWithError if invalid
	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
    cleanedResult := MaskProfanity(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedResult,
		UserID: userID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}
	RespondWithJSON(w, http.StatusCreated, Chirp{
		Body: chirp.Body,
		UserID: chirp.UserID,
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

