package request

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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

	var errs []*ValidationError

	err := validate.Struct(payload)
	if err != nil {
		errSlice := &validator.ValidationErrors{}
		errors.As(err, errSlice)
		for _, err := range *errSlice {
			var element ValidationError
			field, _ := reflect.ValueOf(payload).Type().FieldByName(err.Field())
			element.FailedField = field.Tag.Get("json")
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Type = err.Kind().String()
			element.Message = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", element.FailedField, element.Tag)
			errs = append(errs, &element)
		}
	}
	return errs
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
