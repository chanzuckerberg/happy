package route_groups

import (
	"github.com/chanzuckerberg/happy-api/pkg/cmd/config"
	"github.com/chanzuckerberg/happy-api/pkg/request"
	"github.com/chanzuckerberg/happy-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func RegisterConfig(app *fiber.App) {
	group := app.Group("/config")
	group.Get("/health", request.HealthHandler)

	group.Get("/", func(c *fiber.Ctx) error {
		// TODO: implement
		return c.SendString("list configs")
	})

	group.Get("/:name", func(c *fiber.Ctx) error {
		// TODO: implement
		return c.SendString("get config with name " + c.Params("name"))
	})

	group.Delete("/:name", func(c *fiber.Ctx) error {
		// TODO: implement
		return c.SendString("delete config with name " + c.Params("name"))
	})

	group.Post("/",
		parsePayload[config.SetConfigValuePayload],
		func(c *fiber.Ctx) error {
			payload := c.Context().UserValue("payload").(config.SetConfigValuePayload)

			err := config.SetConfigValue(&payload)
			if err != nil {
				return response.ServerErrorResponse(c, err.Error())
			}

			return c.SendString("set config with name=" + payload.Key + ", value=" + payload.Value)
		},
	)
}

func parsePayload[T interface{}](c *fiber.Ctx) error {
	payload := new(T)
	errors := request.ParsePayload(c, payload)
	if errors != nil {
		return response.ValidationErrorResponse(c, errors)
	}
	c.Context().SetUserValue("payload", *payload)

	return c.Next()
}
