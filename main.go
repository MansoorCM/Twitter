package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()

	rootDir := http.Dir(filePathRoot)
	fileServer := http.FileServer(rootDir)

	mux.Handle("/", fileServer)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("starting server at port 8080")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("error starting server:", err)
	}
}
