package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	const filepathRoot = "."
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("could not retrieve env variable PORT")
	}

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	mux.Handle("/app/*", fsHandler)
	mux.HandleFunc("GET /v1/readiness", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
