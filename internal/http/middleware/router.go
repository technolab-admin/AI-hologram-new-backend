package middleware

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/config"
	"AI-HOLOGRAM-NEW-BACKEND/internal/http/handlers"
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// This file sets up the HTTP and websocket routes

func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	client := meshy.NewClient(cfg.MeshyAPIKey, cfg.MeshyAPIAdress)
	service := meshy.NewService(client)
	handler := handlers.NewMeshyHandler(service)

	r.Route("/meshy", func(r chi.Router) {
		r.Post("/generate", handler.Generate)
	})

	return r
}
