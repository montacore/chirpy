package main

import (
	"net/http"
	"encoding/json"
	"github.com/montacore/chirpy/internal/database"
	"github.com/montacore/chirpy/internal/auth"
	"time"
)


//TEST: "POST /api/login"
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// Decode request body into a parameters struct
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "there was a problem decoding JSON data")
		return
	}


	//TODO: Look up the record by a unique field (not ID)
	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	// TODO: if lookup fails, respond with error, return
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "incorrect email or password, please try again")
		return
	}
	// TODO: Verify the secret -> if error or no match, respond with error, return
	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil  || !match {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	
	


	token, err := auth.MakeJWT(dbUser.ID, cfg.secretKey, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to complete request")
		return
	}
	sessionToken := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: sessionToken,
		UserID: dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60), // 60 days
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error creating refresh token")
		return
	}


	// TODO: map to your response struct, respond 200
	RespondWithJSON(w, http.StatusOK, User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		Token: token,
		RefreshToken: sessionToken, 
		IsChirpyRed: dbUser.IsChirpyRed,
	})
}
