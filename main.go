package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "./"

	mux := http.NewServeMux()

	rootDir := http.Dir(filePathRoot)
	fileServer := http.FileServer(rootDir)

	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("/healthz", handlerHealth)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("starting server at port 8080")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("error starting server:", err)
	}
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
