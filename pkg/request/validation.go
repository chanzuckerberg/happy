package request

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationError struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Message     string `json:"message"`
}

func ValidatePayload(payload interface{}) []*ValidationError {
	validate := validator.New()

	var errors []*ValidationError

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			field, _ := reflect.ValueOf(payload).Type().FieldByName(err.Field())
			element.FailedField = field.Tag.Get("json")
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Type = err.Kind().String()
			element.Message = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", element.FailedField, element.Tag)
			errors = append(errors, &element)
		}
	}
	return errors
}

type ConfigValue struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func ParsePayload[T interface{}](c *fiber.Ctx, payload *T) []*ValidationError {
	if err := c.BodyParser(payload); err != nil {
		ers := []*ValidationError{}
		er := ValidationError{Message: err.Error()}
		ers = append(ers, &er)
		return ers
	}

	return ValidatePayload(*payload)
}
