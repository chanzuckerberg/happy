package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chanzuckerberg/happy-shared/model"
	"github.com/pkg/errors"
)

func ParseResponse[T interface{}](resp *http.Response, result *T) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	return nil
}

func InspectForErrors(resp *http.Response, notFoundMessage string) error {
	if resp.StatusCode == http.StatusNotFound {
		return errors.New(notFoundMessage)
	} else if resp.StatusCode == http.StatusBadRequest {
		validationErrors := []model.ValidationError{}
		ParseResponse(resp, &validationErrors)
		message := ""
		for _, validationError := range validationErrors {
			message = message + fmt.Sprintf("\nhappy-api request failed with: %s", validationError.Message)
		}
		return errors.New(message)
	} else if resp.StatusCode != http.StatusOK {
		errorMessage := new(string)
		ParseResponse(resp, errorMessage)
		return errors.New(*errorMessage)
	}
	return nil
}
