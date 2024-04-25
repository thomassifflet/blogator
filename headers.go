package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/thomassifflet/blogator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiString, err := GetAPIKey(r.Header)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(ctx, apiString)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		handler(w, r, user)
	}
}

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// GetAPIKey -
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed API Key header")
	}

	return splitAuth[1], nil
}
