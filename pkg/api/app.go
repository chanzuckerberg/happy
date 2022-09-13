package api

import (
	"bytes"
	"encoding/json"
	"regexp"

	"github.com/chanzuckerberg/happy-api/pkg/cmd"
	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/request"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
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

type APIApplication struct {
	FiberApp *fiber.App
	DB       *dbutil.DB
}
type APIOption func(*APIApplication)

func WithDebugLogger() APIOption {
	return func(app *APIApplication) {
		app.FiberApp.Use(logger.New(logger.Config{
			Format:     "[${date} ${time}] | ${status} | ${latency} | ${method} | ${path} | ${locals:requestid}\n",
			TimeFormat: "2006-01-02T15:04:05-0700",
			TimeZone:   "UTC",
		}))
	}
}

func WithDatabase(db *dbutil.DB) APIOption {
	return func(app *APIApplication) {
		app.DB = db
	}
}

func MakeAPIApplication() *APIApplication {
	return &APIApplication{
		FiberApp: fiber.New(fiber.Config{
			AppName:     "happy-api",
			JSONEncoder: MarshalJSON,
		}),
		DB: dbutil.MakeDB(),
	}
}

func MakeApp(opts ...APIOption) (*APIApplication, error) {
	apiApp := MakeAPIApplication()
	for _, opt := range opts {
		opt(apiApp)
	}

	err := apiApp.DB.AutoMigrate()
	if err != nil {
		logrus.Fatalf("failed to connect to the DB %s", err)
	}

	apiApp.FiberApp.Use(requestid.New())
	apiApp.FiberApp.Use(func(c *fiber.Ctx) error {
		err := request.VersionCheckHandler(c)
		if err != nil {
			return err
		}
		if c.Response().StatusCode() != fiber.StatusOK {
			return nil
		}
		return c.Next()
	})

	apiApp.FiberApp.Get("/health", request.HealthHandler)
	apiApp.FiberApp.Get("/versionCheck", request.VersionCheckHandler)
	apiApp.FiberApp.Get("/swagger/*", swagger.HandlerDefault)

	v1 := apiApp.FiberApp.Group("/v1")
	RegisterConfigV1(v1, MakeConfigHandler(cmd.MakeConfig(apiApp.DB)))
	RegisterStackListV1(v1, MakeStackHandler(cmd.MakeStack(apiApp.DB)))
	return apiApp, nil
}
