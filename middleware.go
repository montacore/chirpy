package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

