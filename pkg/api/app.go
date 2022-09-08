package api

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"

	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/request"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
)

// Copied from https://gist.github.com/Rican7/39a3dc10c1499384ca91
// with a slight tweak to make "ID" convert to "id" instead of "i_d"
func MarshalJSON(val interface{}) ([]byte, error) {
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	var wordBarrierRegex = regexp.MustCompile(`(\w{2,})([A-Z])`)
	marshalled, err := json.Marshal(val)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return converted, err
}

func MakeApp() (*fiber.App, error) {
	err := dbutil.AutoMigrate()
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiber.Config{
		AppName:     "happy-api",
		JSONEncoder: MarshalJSON,
	})
	app.Use(requestid.New())
	if os.Getenv("APP_ENV") != "test" {
		app.Use(logger.New(logger.Config{
			Format:     "[${date} ${time}] | ${status} | ${latency} | ${method} | ${path} | ${locals:requestid}\n",
			TimeFormat: "2006-01-02T15:04:05-0700",
			TimeZone:   "UTC",
		}))
	}
	app.Use(func(c *fiber.Ctx) error {
		err := request.VersionCheckHandler(c)
		if err != nil {
			return err
		}
		if c.Response().StatusCode() != fiber.StatusOK {
			return nil
		}
		return c.Next()
	})

	app.Get("/health", request.HealthHandler)
	app.Get("/versionCheck", request.VersionCheckHandler)
	app.Get("/swagger/*", swagger.HandlerDefault)

	v1 := app.Group("/v1")
	RegisterConfigV1(&v1)

	return app, nil
}
