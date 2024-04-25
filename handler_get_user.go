package main

import (
	"context"
	"net/http"

	"github.com/thomassifflet/blogator/internal/database"
)

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	ctx := context.Background()
	type response struct {
		User
	}

	authString, err := GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 500, "no api key provided")
		return
	}

	userByKey, err := cfg.DB.GetUserByAPIKey(ctx, authString)
	if err != nil {
		respondWithError(w, 500, "couldn't retrieve user")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        userByKey.ID,
			CreatedAt: userByKey.CreatedAt,
			UpdatedAt: userByKey.CreatedAt,
			Name:      userByKey.Name,
			ApiKey:    userByKey.ApiKey,
		},
	})

}
