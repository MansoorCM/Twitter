package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	const port = "8080"
	const filePathRoot = "./"
	apiCfg := apiConfig{fileServerHits: atomic.Int32{}}

	mux := http.NewServeMux()

	rootDir := http.Dir(filePathRoot)
	fileServer := http.FileServer(rootDir)

	mux.Handle("/app/", apiCfg.middleWareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("starting server at port 8080")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("error starting server:", err)
	}
}
