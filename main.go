package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(FiberConfig)
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", Running)
	app.Get("/api/tags", ModelList)
	app.Post("/api/chat", ChatCompletion)

	err := app.Listen(":11434")
	if err != nil {
		log.Fatal(err)
	}
}
