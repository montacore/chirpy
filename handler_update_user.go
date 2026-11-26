package main

import (
	"encoding/json"
	"net/http"
	"github.com/montacore/chirpy/internal/auth"
	"github.com/montacore/chirpy/internal/database"
)


func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "unable to validate")
		return
	}
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "unable to decode json")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "unable to update email or password")
		return
	}

	RespondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})

}
