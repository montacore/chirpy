package main

import (
	"net/http"
	"github.com/montacore/chirpy/internal/auth"
	"time"
)
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "unable to retrieve bearer token")
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "unable to retrieve user")
		return
	}
	newToken, err := auth.MakeJWT(user.ID,  cfg.secretKey, time.Hour)
	if err != nil {
		
		RespondWithError(w, http.StatusInternalServerError, "unable to create token")
		return
	}

	RespondWithJSON(w, http.StatusOK, response {
		Token: newToken,
	})

}
