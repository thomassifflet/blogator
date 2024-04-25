package main

import "net/http"

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feedsDB, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt retrieve feeds")
		return
	}

	respondWithJSON(w, 200, feedsDB)

}
