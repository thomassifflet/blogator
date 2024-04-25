package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	srv := &http.Server{
		Addr: ":" + port,
	}

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	r.Handle("/app/*", fsHandler)
	r.HandleFunc("/v1/users", apiCfg.handlerCreateUser).Methods("POST")
	r.HandleFunc("/v1/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser)).Methods("GET")
	r.HandleFunc("/v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed)).Methods("POST")
	r.HandleFunc("/v1/readiness", handlerReadiness).Methods("GET")
	r.HandleFunc("/v1/err", handlerErr).Methods("GET")

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
