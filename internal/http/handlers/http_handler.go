package handlers

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/config"
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type HttpHandler struct {
	meshy *meshy.Service
	cfg   *config.Config
}

func NewHttpHandler(meshyService *meshy.Service, c *config.Config) *HttpHandler {
	return &HttpHandler{meshy: meshyService, cfg: c}
}

func (h *HttpHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req meshy.TextTo3DRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()

	go h.runMeshyJob(jobID, &req)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

func (h *HttpHandler) runMeshyJob(jobID string, req *meshy.TextTo3DRequest) {

}
