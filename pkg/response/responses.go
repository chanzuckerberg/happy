package response

import (
	"github.com/gofiber/fiber/v2"
)

func ValidationErrorResponse(c *fiber.Ctx, errors []*ValidationError) error {
	return c.Status(fiber.StatusBadRequest).JSON(errors)
}

func ServerErrorResponse(c *fiber.Ctx, errorMessage string) error {
	serverErr := map[string]string{"message": errorMessage}
	errors := [1]map[string]string{serverErr}
	return c.Status(fiber.StatusInternalServerError).JSON(errors)
}
