package handlers

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/logger"
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type JobRunner interface {
	Run(id string, req *meshy.TextTo3DRequest)
}

type HttpHandler struct {
	jobRunner JobRunner
}

func NewHttpHandler(jobRunner JobRunner) *HttpHandler {
	return &HttpHandler{jobRunner: jobRunner}
}

func (h *HttpHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req meshy.TextTo3DRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn.Printf("invalid JSON body: %v", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()

	go h.jobRunner.Run(jobID, &req)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})

}
