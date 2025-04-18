package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ChatCompletion(c *fiber.Ctx) error {
	// Parse ollama request
	chat, err := BindOllamaChat(c.BodyParser)
	if err != nil {
		log.Error(err)
		return err
	}

	// Get the appropriate provider for the model
	provider := GetProvider(chat.Model)

	// Use the provider to handle the chat completion
	return provider.ChatCompletion(c, chat)
}
