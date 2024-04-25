package main

import "net/http"

type response struct {
	Resp string `json:"status"`
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	resp := response{
		Resp: "ok",
	}
	respondWithJSON(w, 200, resp)
}
