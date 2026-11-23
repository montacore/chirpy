package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"github.com/montacore/chirpy/internal/database"
	"github.com/montacore/chirpy/internal/auth"
)
//NOTE: the `atomic.Int32` is a standard library type that allows us to
// safely increment and read an integer value across multiple goroutines(HTTP reqs)
type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

// NOTE: middlewareMetricsInc takes one argument:
// 	next 		http.Handler
// 	and returns an instance of:
// 				http.Handler
// The function returns the result of a call to HandlerFunc(func(w http.ResponseWriter, r *http.Request))
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
	
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(
		w, 
		`
<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
</html>
		`,
		cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if cfg.platform != "dev" {
		RespondWithError(w, http.StatusForbidden, "You are not authorized to reset database")
		return
	}
	err  := cfg.db.Reset(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error resetting database")
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset complete"))
}


func (cfg *apiConfig)handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body 	string		`json:"body"`
		UserID	uuid.UUID		`json:"user_id"`
	}

	// TODO: Decode r.Body into parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		
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
		UserID: params.UserID,
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

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type parameters struct {
		Password string `json:"hashed_password"`
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
		RespondWithError(w, http.StatusBadRequest, "error hashing password")
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
		
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chirps, err := cfg.db.GetChirps(ctx)
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
