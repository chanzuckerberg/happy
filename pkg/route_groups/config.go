package route_groups

import (
	"github.com/chanzuckerberg/happy-api/pkg/cmd/config"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/chanzuckerberg/happy-api/pkg/request"
	"github.com/chanzuckerberg/happy-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func RegisterConfig(app *fiber.App) {
	group := app.Group("/config")
	group.Get("/health", request.HealthHandler)

	// debugging endpoint that returns all config values for an app+env combo without resolving
	group.Get("/dump", parsePayload[model.AppMetadata], func(c *fiber.Ctx) error {
		payload := getPayload[model.AppMetadata](c)
		records := config.GetAllAppConfigs(&payload)

		return c.Status(fiber.StatusOK).JSON(wrapWithCount(&records))
	})

	loadConfigs(app)
}

func loadConfigs(app *fiber.App) {
	group := app.Group("/configs")

	group.Get("/", parsePayload[model.AppMetadata], func(c *fiber.Ctx) error {
		payload := getPayload[model.AppMetadata](c)

		var records []*model.AppConfigResponse
		if payload.Stack == "" {
			records = config.GetAppConfigsForEnv(&payload)
		} else {
			records = config.GetAppConfigsForStack(&payload)
		}

		return c.Status(fiber.StatusOK).JSON(wrapWithCount(&records))
	})

	group.Post("/",
		parsePayload[model.AppConfigPayload],
		func(c *fiber.Ctx) error {
			payload := getPayload[model.AppConfigPayload](c)
			config.SetConfigValue(&payload)

			return c.Status(fiber.StatusOK).JSON(map[string]string{
				"key":   payload.Key,
				"value": payload.Value,
			})
		})

	group.Get("/:key",
		parsePayload[model.AppMetadata],
		func(c *fiber.Ctx) error {
			payload := model.AppConfigLookupPayload{
				AppMetadata: getPayload[model.AppMetadata](c),
				ConfigKey:   model.ConfigKey{Key: c.Params("key")},
			}
			record := config.GetAppConfig(&payload)

			status := c.Status(fiber.StatusOK)
			if record == nil {
				status = c.Status(fiber.StatusNotFound)
			}

			return status.JSON(map[string]interface{}{"record": record})
		})

	group.Delete("/:key",
		parsePayload[model.AppMetadata],
		func(c *fiber.Ctx) error {
			payload := model.AppConfigLookupPayload{
				AppMetadata: getPayload[model.AppMetadata](c),
				ConfigKey:   model.ConfigKey{Key: c.Params("key")},
			}
			record := config.DeleteAppConfig(&payload)

			return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
				"deleted": record == nil,
				"record":  record,
			})
		})
}

func getPayload[T interface{}](c *fiber.Ctx) T {
	return c.Context().UserValue("payload").(T)
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

func wrapWithCount[T interface{}](records *[]*T) *map[string]interface{} {
	return &map[string]interface{}{
		"records": records,
		"count":   len(*records),
	}
}
