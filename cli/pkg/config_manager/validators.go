package config_manager

import (
	"reflect"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

func PortValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return errors.Errorf("cannot enforce a port number check on response of type %v", reflect.TypeOf(val).Name())
	}
	if !govalidator.IsPort(str) {
		return errors.New("value is not a valid port number")
	}
	return nil
}

func URIValidator(val interface{}) error {
	str, ok := val.(string)

	if !ok {
		return errors.Errorf("cannot enforce a uri check on response of type %v", reflect.TypeOf(val).Name())
	}

	if !govalidator.IsRequestURI(str) {
		return errors.New("value is not a valid uri")
	}

	return nil
}
