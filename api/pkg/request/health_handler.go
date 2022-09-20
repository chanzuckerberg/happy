package request

import (
	"github.com/gofiber/fiber/v2"
)

func HealthHandler(c *fiber.Ctx) error {
	status := map[string]string{"status": "OK", "route": c.Path()}
	return c.JSON(status)
}
