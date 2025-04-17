package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func Running(c *fiber.Ctx) error {
	return c.SendString("Ollama is running")
}

func ModelList(c *fiber.Ctx) error {
	pages, err := provider.Models.List(context.Background())
	if err != nil {
		return err
	}
	tags := Tags{
		Models: []Model{},
	}
	for _, model := range pages.Data {
		tags.Models = append(tags.Models, Model{
			Name:  model.ID,
			Model: model.ID,
		})
	}
	return c.JSON(tags)
}
