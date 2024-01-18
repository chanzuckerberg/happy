package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Response struct {
	Status   string
	Service  string
	Env      string
	Complete bool
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
		env := strings.Join(os.Environ(), "\n")
		return c.Status(http.StatusOK).JSON(Response{Status: "OK", Service: "frontend", Env: env})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		env := strings.Join(os.Environ(), "\n")
		return c.Status(http.StatusOK).JSON(Response{Status: "Health", Service: "frontend", Env: env})
	})

	app.Get("/proxy", func(c *fiber.Ctx) error {
		internalApiEndpoint := os.Getenv("PRIVATE_INTERNAL_API_ENDPOINT")
		count := 1
		count, err := strconv.Atoi(c.Query("count"))
		if err != nil {
			count = 1
		}

		body := []byte{}
		for i := 0; i < count; i++ {
			resp, err := http.Get(fmt.Sprintf("%s/", internalApiEndpoint))
			if err != nil {
				return c.Status(http.StatusOK).JSON(Response{Status: fmt.Sprintf("Error making http call: %s", err.Error()), Service: "frontend", Complete: false})
			}
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return c.Status(http.StatusOK).JSON(Response{Status: fmt.Sprintf("Error reading response: %s", err.Error()), Service: "frontend", Complete: false})
			}
		}
		return c.Status(http.StatusOK).JSON(Response{Status: string(body), Service: "frontend", Complete: true})
	})

	app.Listen(":3000")
}
