package auth

import (
	"net/http"
	"errors"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authorization header")
	}
	key := strings.TrimPrefix(val, "ApiKey ")
	return strings.TrimSpace(key), nil
}
