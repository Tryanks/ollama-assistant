package main

import (
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default or system environment variables")
	}

	// Scan models at initialization
	if err := ScanModels(); err != nil {
		fiberlog.Warn("Failed to scan models at initialization:", err)
	}
}

var FiberConfig = fiber.Config{
	JSONEncoder: sonic.Marshal,
	JSONDecoder: sonic.Unmarshal,
}

// getEnv retrieves the value of the environment variable named by the key
// If the variable is not present, it returns the fallback value
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

// getEnvWithDefault retrieves the value of the environment variable named by the key
// If the variable is not present, it returns the provided default value
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
