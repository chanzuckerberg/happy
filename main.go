package main

import (
	"github.com/chanzuckerberg/happy-api/pkg/handlers"
	"github.com/chanzuckerberg/happy-api/pkg/route_groups"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "happy-api",
	})

	app.Get("/health", handlers.Status)

	route_groups.RegisterConfig(app)

	app.Listen(":3001")
}
