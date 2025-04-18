package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"slices"
	"strings"
)

func Running(c *fiber.Ctx) error {
	return c.SendString("Ollama is running")
}

func ModelList(c *fiber.Ctx) error {
	pages, err := provider.Models.List(context.Background())
	if err != nil {
		return err
	}
	uniques := make(map[string]struct{})
	for _, model := range pages.Data {
		if ModelBlockFilter.BlockString(model.ID) {
			continue
		}
		uniques[model.ID] = struct{}{}
	}
	tags := Tags{
		Models: []Model{},
	}
	for model := range uniques {
		tags.Models = append(tags.Models, Model{
			Name:  model,
			Model: model,
		})
	}
	slices.SortFunc(tags.Models, func(a, b Model) int {
		return strings.Compare(a.Name, b.Name)
	})
	return c.JSON(tags)
}
