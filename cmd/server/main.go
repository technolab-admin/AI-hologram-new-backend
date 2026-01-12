package main

// This is where the server gets started up

import (
	"fmt"
	"net/http"

	"AI-HOLOGRAM-NEW-BACKEND/internal/config"
	"AI-HOLOGRAM-NEW-BACKEND/internal/http/middleware"
	"AI-HOLOGRAM-NEW-BACKEND/internal/logger"
)

func main() {
	logger.Init()
	if err := run(); err != nil {
		logger.Error.Fatal(err)
	}
}

func run() error {
	logger.Info.Println("Server starting...")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}

	r := middleware.NewRouter(cfg)

	logger.Info.Printf("Server running on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}
