package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
fileserverHits atomic.Int32
}

var cfg apiConfig

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
