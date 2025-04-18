package main

import (
	"github.com/gofiber/fiber/v2"
)

func Running(c *fiber.Ctx) error {
	return c.SendString("Ollama is running")
}

func ModelList(c *fiber.Ctx) error {
	// Return the cached model list
	modelListMutex.RLock()
	defer modelListMutex.RUnlock()
	return c.JSON(cachedModelList)
}

// RefreshModels refreshes the model list and returns the updated list
func RefreshModels(c *fiber.Ctx) error {
	// Scan models and update the cached list
	err := ScanModels()
	if err != nil {
		return err
	}

	// Return the updated model list
	modelListMutex.RLock()
	defer modelListMutex.RUnlock()
	return c.JSON(cachedModelList)
}
