package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middleWareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	htmlResponse := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerHits.Load())
	fmt.Fprint(w, htmlResponse)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithJson(w, []byte("Reset is only allowed in Dev environment."), http.StatusForbidden)
	}
	cfg.fileServerHits.Store(0)
	cfg.db.Reset(r.Context())
	respondWithJson(w, []byte("Hits reset to 0."), http.StatusOK)

}
