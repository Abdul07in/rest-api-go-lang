package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
