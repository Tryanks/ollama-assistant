package main

import (
	"github.com/openai/openai-go/option"
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
)

func init() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default or system environment variables")
	}

	provider = openai.NewClient(
		option.WithBaseURL(getEnv("API_BASE_URL")),
		option.WithAPIKey(getEnv("API_KEY")),
	)
}

var FiberConfig = fiber.Config{
	JSONEncoder: sonic.Marshal,
	JSONDecoder: sonic.Unmarshal,
}

var provider openai.Client

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
