package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chanzuckerberg/happy-shared/model"
	"github.com/pkg/errors"
)

var ErrRecordNotFound = errors.New("record not found")

func ParseResponse[T interface{}](resp *http.Response, result *T) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	return nil
}

func InspectForErrors(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrRecordNotFound
	case http.StatusBadRequest:
		validationErrors := []model.ValidationError{}
		ParseResponse(resp, &validationErrors)
		message := ""
		for _, validationError := range validationErrors {
			message = message + fmt.Sprintf("\nhappy-api request failed with: %s", validationError.Message)
		}
		return errors.New(message)
	default:
		errorMessage := new(string)
		ParseResponse(resp, errorMessage)
		return errors.New(*errorMessage)
	}
}
