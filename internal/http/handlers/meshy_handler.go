package handlers

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/config"
	"AI-HOLOGRAM-NEW-BACKEND/internal/logger"
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
	"encoding/json"
	"fmt"
	"net/http"
)

type MeshyHandler struct {
	service *meshy.Service
	cfg     *config.Config
}

func NewMeshyHandler(s *meshy.Service, c *config.Config) *MeshyHandler {
	return &MeshyHandler{service: s, cfg: c}
}

func (h *MeshyHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req meshy.TextTo3DRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn.Printf("invalid JSON body: %v", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	previewID, err := h.service.GeneratePreview(&req)

	if err != nil {
		logger.Error.Printf("preview generation failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refineName, err := h.service.GenerateRefine(previewID)

	if err != nil {
		logger.Error.Printf("refine generation failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info.Printf("generation succeeded: model_url=%s", refineName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"model_url": fmt.Sprintf("%s/assets/downloads/%s", h.cfg.PublicBaseUrl, refineName),
	})
}
