package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/montacore/chirpy/internal/auth"
)

func (cfg *apiConfig)handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "unable to authenticate")
		return
	}
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "unable to locate chirp")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "unable to authenticate token")
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		//TODO: If chirp not found return 404 error
		RespondWithError(w, http.StatusNotFound, "item not found")
		return
	}
	if chirp.UserID != userID {
		RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "there was a problem deleting the chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)


}
