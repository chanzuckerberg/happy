package client

import (
	"encoding/json"
	"net/http"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/hashicorp/go-multierror"
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
		var errs error
		for _, err := range validationErrors {
			errs = multierror.Append(errs, err)
		}
		return errs
	default:
		errorMessage := new(string)
		ParseResponse(resp, errorMessage)
		return errors.New(*errorMessage)
	}
}
