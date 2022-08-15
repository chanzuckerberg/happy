package main

import (
	"github.com/chanzuckerberg/happy-api/pkg/request"
	"github.com/chanzuckerberg/happy-api/pkg/route_groups"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

func exec() error {
	app := fiber.New(fiber.Config{
		AppName: "happy-api",
	})
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${date} ${time}] | ${status} | ${latency} | ${method} | ${path} | ${locals:requestid}\n",
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "UTC",
	}))

	app.Get("/health", request.HealthHandler)

	route_groups.RegisterConfig(app)

	return app.Listen(":3001")
}

func main() {
	err := exec()
	if err != nil {
		logrus.Error(err)
	}
}
