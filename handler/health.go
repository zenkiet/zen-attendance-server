package handler

import (
	"net/http"

	"zenkiet/zen-attendance-server/pkg/response"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func HandleHealth(w *http.ResponseWriter, r *http.Request) {
	response.JSON(*w, http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "zen-attendance-server",
		Version: "1.0.0",
	})
}
