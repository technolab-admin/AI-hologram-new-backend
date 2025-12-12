package handlers

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
	"encoding/json"
	"fmt"
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
		fmt.Println("JSON Decode error:", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("REQUEST: %+v\n", req) // debug

	url, err := h.service.Generate(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"model_url": url,
	})
}
