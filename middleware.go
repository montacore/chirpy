package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"github.com/montacore/chirpy/internal/database"
)

//NOTE: the `atomic.Int32` is a standard library type that allows us to
// safely increment and read an integer value across multiple goroutines(HTTP reqs)
type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secretKey string
}

// NOTE: middlewareMetricsInc takes one argument:
// 	next 		http.Handler
// 	and returns an instance of:
// 				http.Handler
// The function returns the result of a call to HandlerFunc(func(w http.ResponseWriter, r *http.Request))
//TEST: entry endpoint "/app/"
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
	
}

//TEST: Endpoint: "GET /admin/metrics"
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(
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

//TEST: Endpoint: "POST /admin/reset"
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
	_,_ = w.Write([]byte("Reset complete"))
}

