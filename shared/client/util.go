package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrUnauthorized   = errors.New("unauthorized")
)

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
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrRecordNotFound
	case http.StatusBadRequest:
		validationErrors := []model.ValidationError{}
		err := ParseResponse(resp, &validationErrors)
		if err != nil {
			return errors.Wrapf(err, "unable to parse resp body as JSON for status code %+v", http.StatusBadRequest)
		}
		var errs error
		for _, err := range validationErrors {
			errs = multierror.Append(errs, err)
		}
		return errs
	default:
		var errorMessage interface{}
		err := ParseResponse(resp, &errorMessage)
		if err != nil {
			return errors.Wrapf(err, "unable to parse resp body as JSON for status code %+v", resp.StatusCode)
		}
		return errors.New(fmt.Sprintf("status code %+v: %+v", resp.StatusCode, errorMessage))
	}
}
