package auth


import (
	"github.com/alexedwards/argon2id"
	"log"
	"errors"
	"strings"
	"net/http"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("encountered an error hashing password: %s\n", err)
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("encountered error when comparing password hashes: %s\n", err)
		return false, err
	}
	return match, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authorization header")
	}
	key := strings.TrimPrefix(val, "Bearer ")
	return strings.TrimSpace(key), nil
}
