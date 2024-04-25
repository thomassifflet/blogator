package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thomassifflet/blogator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		User
	}

	authString, err := GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 500, "no api key provided")
		return
	}

	userByKey, err := cfg.DB.GetUserByAPIKey(r.Context(), authString)
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

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	type parameters struct {
		Name string `json:"name"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	userDb, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        userDb.ID,
			CreatedAt: userDb.CreatedAt,
			UpdatedAt: userDb.UpdatedAt,
			Name:      userDb.Name,
		},
	})
}
