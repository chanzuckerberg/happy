package route_groups

import (
	"github.com/chanzuckerberg/happy-api/pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterConfig(app *fiber.App) {
	group := app.Group("/config")
	group.Get("/health", handlers.Status)

}
