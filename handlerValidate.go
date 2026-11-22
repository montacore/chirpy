package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat,  err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}


func RespondWithError(n http.ResponseWriter, code int, msg string){
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(n, code, errorResponse{Error: msg})
}

func MaskProfanity(text string) string {
	banned := map[string]bool{
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	words := strings.Split(text, " ")
	for i, w := range words {
		if banned[strings.ToLower(w)] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
