package main

import (
	"net/http"
	"github.com/montacore/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header) 
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "unable to retrieve bearer token")
		return
	}
	err = cfg.db.RevokeToken(r.Context(), token)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to revoke token")
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}

