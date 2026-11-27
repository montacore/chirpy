package main

import (
	"net/http"
	"encoding/json"
	"github.com/montacore/chirpy/internal/database"
	"github.com/montacore/chirpy/internal/auth"
)


// TEST: Endpoint: "POST /api/users"
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "there was a problem decoding JSON data")
		return
	}
	pass, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	dbUser, err := cfg.db.CreateUser(ctx, database.CreateUserParams{
		Email: params.Email,
		HashedPassword: pass,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "there was a problem creating a new user")
		return
	}
	RespondWithJSON(w, http.StatusCreated, User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	})

}
