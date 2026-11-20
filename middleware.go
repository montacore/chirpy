package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"encoding/json"
	"log"
)
//NOTE: the `atomic.Int32` is a standard library type that allows us to
// safely increment and read an integer value across multiple goroutines(HTTP reqs)
type apiConfig struct {
	fileserverHits atomic.Int32
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
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}


func (cfg *apiConfig)handlerValidateBody(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body 	string		`json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
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
	RespondWithJSON(w, http.StatusOK, response{
		CleanedBody: cleanedResult,
	})
	
}
