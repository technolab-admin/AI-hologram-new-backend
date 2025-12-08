package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var loadOnce sync.Once

func Load() {
	loadOnce.Do(func() {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Printf("%s", err)
			log.Println("No .env file found (proceeding with system environment)")
		}
	})
}

func GetString(key string) string {
	Load()
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

func GetInt(key string) int {
	Load()
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}

	i, err := strconv.Atoi(val)

	if err != nil {
		log.Fatalf("Invalid integer value for environment variable %s: %v", key, err)
	}

	return i
}

func GetOptionalString(key string, fallback string) string {
	Load()
	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	return val
}
