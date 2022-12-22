package request

import (
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/gofiber/fiber/v2"
)

// HealthCheck godoc
// @Summary Check that the server is up.
// @Tags    root
// @Accept  */*
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router  /health [get]
func HealthHandler(c *fiber.Ctx) error {
	status := model.HealthResponse{Status: "OK", Route: c.Path()}
	return c.JSON(status)
}
