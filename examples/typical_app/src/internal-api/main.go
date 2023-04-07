package main

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Response struct {
	Status  string
	Service string
}

func main() {
	app := fiber.New(fiber.Config{
		ReadTimeout:    60 * time.Second,
		ReadBufferSize: 1024 * 64,
	})
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path} ${reqHeaders}​\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(Response{Status: "OK", Service: "internal-api"})
	})

	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(Response{Status: "OK", Service: "internal-api"})
	})

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(Response{Status: "Health", Service: "internal-api"})
	})

	app.Listen(":3000")
}
