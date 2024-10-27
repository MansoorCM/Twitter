package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()

	dir := http.Dir(filePathRoot)
	fileServer := http.FileServer(dir)

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
