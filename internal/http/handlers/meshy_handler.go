package handlers

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/logger"
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
	"encoding/json"
	"net/http"
)

type MeshyHandler struct {
	service *meshy.Service
}

func NewMeshyHandler(s *meshy.Service) *MeshyHandler {
	return &MeshyHandler{service: s}
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

	refineUrl, err := h.service.GenerateRefine(previewID)

	if err != nil {
		logger.Error.Printf("refine generation failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info.Printf("generation succeeded: model_url=%s", refineUrl)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"model_url": refineUrl,
	})
}
