package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr     string
	WebsocketAddr  string
	WebsocketURL   string
	MeshyAPIKey    string
	MeshyAPIAdress string
	PublicBaseUrl  string
}

func Load() (*Config, error) {
	godotenv.Load(".env")

	return &Config{
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		WebsocketAddr:  getEnv("WEBSOCKET_ADDR", ":8081"),
		WebsocketURL:   os.Getenv("WEBSOCKET_URL"),
		MeshyAPIKey:    os.Getenv("MESHY_API_KEY"),
		MeshyAPIAdress: os.Getenv("MESHY_API_BASE_URL"),
		PublicBaseUrl:  os.Getenv("PUBLIC_BASE_URL"),
	}, nil
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
