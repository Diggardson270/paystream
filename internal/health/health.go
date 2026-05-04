package health

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{
		Status:  "ok",
		Service: "paystream-api",
		Version: "0.1.0",
	})
}
