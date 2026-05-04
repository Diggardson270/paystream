package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/breedar/paystream/internal/health"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.Handler)

	log.Println("paystream-api listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// notFound is the default 404 handler used until the full router lands.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
}
