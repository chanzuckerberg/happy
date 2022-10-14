package request

import (
	"github.com/chanzuckerberg/happy/api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// HealthCheck godoc
// @Summary Check that the server is up.
// @Tags    root
// @Accept  */*
// @Produce json
// @Success 200 {object} response.HealthResponse
// @Router  /health [get]
func HealthHandler(c *fiber.Ctx) error {
	status := response.HealthResponse{Status: "OK", Route: c.Path()}
	return c.JSON(status)
}
