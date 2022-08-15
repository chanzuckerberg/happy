package config

type AppMetadata struct {
	AppName     string `json:"app" validate:"required"`
	Environment string `json:"environment" validate:"required"`
	Stack       string `json:"stack"`
}

type ConfigValue struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type SetConfigValuePayload struct {
	AppMetadata
	ConfigValue
}

func SetConfigValue(payload *SetConfigValuePayload) error {
	// TODO: implement
	return nil
}
