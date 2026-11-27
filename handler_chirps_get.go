package main

import (
	"net/http"
	"log"
	"github.com/google/uuid"
	"github.com/montacore/chirpy/internal/database"
) 
//TEST: Endpoint: "GET /api/chirps"
func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var chirps []database.Chirp
	var err error
	s := r.URL.Query().Get("author_id")
	if s == "" {
	chirps, err = cfg.db.GetChirps(ctx)
	} else {
		id, parseErr := uuid.Parse(s)
		if parseErr != nil {
			RespondWithError(w, http.StatusBadRequest, "unable to parse author_id")
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(ctx, id)
	}
	if err != nil {
		log.Printf("error getting chirps: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "there was a problem retrieving chirps")
		return
	}
	out := []Chirp{}
	for _, o := range chirps {
		out = append(out, Chirp{
			ID: o.ID,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			Body: o.Body,
			UserID: o.UserID,
		})
	}
	RespondWithJSON(w, http.StatusOK, out)
}


//TEST: "GET /api/chirps/{chirpID}"
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("[error] unable to parse the requested chirp id %s\n", err)
		RespondWithError(w, http.StatusBadRequest, "unable to parse the requested id string")
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		log.Printf("[error] the requested chirp was not found: %s\n", err)
		RespondWithError(w, http.StatusNotFound, "the requested chirp was not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID: chirp.UserID,
		Body: chirp.Body,
	})
}
