package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Status(c *fiber.Ctx) error {
	status := map[string]string{"status": "OK", "route": c.Path()}
	return c.JSON(status)
}
