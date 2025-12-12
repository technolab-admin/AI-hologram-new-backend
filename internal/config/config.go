package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr  string
	MeshyAPIKey string
}

func Load() (*Config, error) {
	godotenv.Load(".env")

	return &Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		MeshyAPIKey: os.Getenv("MESHY_API_KEY"),
	}, nil
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
