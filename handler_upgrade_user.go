package main

import (
	"github.com/google/uuid"
	"net/http"
	"encoding/json"
	"database/sql"
	"errors"
	"github.com/montacore/chirpy/internal/auth"
)

func (cfg *apiConfig)handlerVerifyUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "couldn't find api key")
		return
	}
	if key != cfg.polkaKey {
		RespondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "couldn't decode parameters into parameters")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = cfg.db.VerifyUpgrade(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, http.StatusNotFound, "unable to locate upgrade")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "couldn't locate upgrade")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
