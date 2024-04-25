package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/thomassifflet/blogator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	const filepathRoot = "."
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("could not retrieve env variable PORT")
	}

	dbURL := os.Getenv("DBCONN")
	if dbURL == "" {
		log.Fatal("could not retrieve database connection URL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
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
	mux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUser)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
