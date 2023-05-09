package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Response struct {
	Status   string
	Service  string
	Complete bool
}

const sidecarEndpoint = "http://localhost"

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
		return c.Status(http.StatusOK).JSON(Response{Status: "OK", Service: "frontend"})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(Response{Status: "Health", Service: "frontend"})
	})

	app.Get("/sidecar", func(c *fiber.Ctx) error {
		body := []byte{}

		resp, err := http.Get(sidecarEndpoint)
		if err != nil {
			return c.Status(http.StatusOK).JSON(Response{Status: fmt.Sprintf("Error making http call: %s", err.Error()), Service: "frontend", Complete: false})
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(http.StatusOK).JSON(Response{Status: fmt.Sprintf("Error reading response: %s", err.Error()), Service: "frontend", Complete: false})
		}

		return c.Status(http.StatusOK).JSON(Response{Status: string(body), Service: "frontend", Complete: true})
	})

	app.Listen(":3000")
}
