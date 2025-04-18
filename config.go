package main

import (
	"log"
	"os"
	"strings"

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

// OpenAIProvider represents an OpenAI API provider with its base URL and API key
type OpenAIProvider struct {
	BaseURL string
	APIKey  string
	Models  []string
}

// GetOpenAIProviders parses the OPENAI_PROVIDERS environment variable and returns a slice of OpenAIProvider
// Format: OPENAI_PROVIDERS=url1, key1; url2, key2;...
func GetOpenAIProviders() []OpenAIProvider {
	providersStr := getEnv("OPENAI_PROVIDERS")

	// If OPENAI_PROVIDERS is not set, use the legacy API_BASE_URL and API_KEY
	if providersStr == "" {
		return []OpenAIProvider{
			{
				BaseURL: getEnv("API_BASE_URL"),
				APIKey:  getEnv("API_KEY"),
				Models:  []string{},
			},
		}
	}

	var providers []OpenAIProvider

	// Split by semicolon to get each provider
	providerParts := strings.Split(providersStr, ";")
	for _, part := range providerParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split by comma to get URL and key
		urlAndKey := strings.SplitN(part, ",", 2)
		if len(urlAndKey) != 2 {
			fiberlog.Warn("Invalid provider format:", part)
			continue
		}

		providers = append(providers, OpenAIProvider{
			BaseURL: strings.TrimSpace(urlAndKey[0]),
			APIKey:  strings.TrimSpace(urlAndKey[1]),
			Models:  []string{},
		})
	}

	// If no valid providers were found, fall back to legacy configuration
	if len(providers) == 0 {
		return []OpenAIProvider{
			{
				BaseURL: getEnv("API_BASE_URL"),
				APIKey:  getEnv("API_KEY"),
				Models:  []string{},
			},
		}
	}

	return providers
}
