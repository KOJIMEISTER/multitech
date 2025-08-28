package config

import (
	"log"
	"os"
)

func LoadEnv() {
	required := []string{
		"JWT_SECRET",
		"REDIS_URL",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required environment variable: %s", key)
		}
	}
}
