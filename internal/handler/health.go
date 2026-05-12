package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse is the JSON body for GET /health.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health returns a minimal liveness probe (no DB check).
func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	}
}
