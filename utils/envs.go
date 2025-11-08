package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// load envs from .env file
func LoadEnvs() {
	// Try multiple possible locations for .env file
	paths := []string{
		".env",       // current directory
		"../.env",    // parent directory
		"../../.env", // grandparent directory
	}

	// Also try to find .env from the module root
	if wd, err := os.Getwd(); err == nil {
		// Walk up the directory tree looking for go.mod
		dir := wd
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				paths = append(paths, filepath.Join(dir, ".env"))
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// Try each path until one works
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			return
		}
	}

	log.Println("Warning: Could not load .env file from any common location")
}

func GetApiKeyFromEnv() (string, bool) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", false
	}
	return apiKey, true
}
