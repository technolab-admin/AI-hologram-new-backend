package main

// This is where the server gets started up

import (
	"log"
	"net/http"
	
	"AI-HOLOGRAM-NEW-BACKEND/internal/websockets"
	"AI-HOLOGRAM-NEW-BACKEND/internal/config"
	"AI-HOLOGRAM-NEW-BACKEND/internal/http/middleware"
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	wsServer := websockets.NewServer(cfg.WebsocketAddr)
	go wsServer.Start()


	wsClient := meshy.NewWSClient("backend-meshy")
	r := middleware.NewRouter(cfg, wsClient)
	go wsClient.StartWebsocketClient()


	log.Printf("Server running on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	
}
