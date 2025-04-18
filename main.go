package main

import (
	"fmt"
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

	err := app.Listen(fmt.Sprintf("%s:%s",
		getEnvWithDefault("HOST_SERVE", ""),
		getEnvWithDefault("PORT_SERVE", "11434"),
	))
	if err != nil {
		log.Fatal(err)
	}
}
