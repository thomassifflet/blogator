package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	r.HandleFunc("/v1/users", apiCfg.handlerCreateUser).Methods("POST")
	r.HandleFunc("/v1/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser)).Methods("GET")
	r.HandleFunc("/v1/feeds", apiCfg.handlerGetFeeds).Methods("GET")
	r.HandleFunc("/v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed)).Methods("POST")
	r.HandleFunc("/v1/readiness", handlerReadiness).Methods("GET")
	r.HandleFunc("/v1/err", handlerErr).Methods("GET")
	r.HandleFunc("/v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsCreate)).Methods("POST")
	r.HandleFunc("/v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet)).Methods("GET")
	r.HandleFunc("/v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete)).Methods("DELETE")

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(http.ListenAndServe(port, r))
}
