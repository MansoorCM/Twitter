package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MansoorCM/Twitter/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set.")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set.")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET must be set.")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("couldn't open db %s", err)
	}

	dbQueries := database.New(db)

	const port = "8080"
	const filePathRoot = "./"
	apiCfg := apiConfig{fileServerHits: atomic.Int32{}, db: dbQueries, platform: platform, jwtSecret: secret}

	mux := http.NewServeMux()

	rootDir := http.Dir(filePathRoot)
	fileServer := http.FileServer(rootDir)

	mux.Handle("/app/", apiCfg.middleWareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirp)

	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("POST /api/login", apiCfg.userLogin)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("starting server at port 8080")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("error starting server:", err)
	}
}
