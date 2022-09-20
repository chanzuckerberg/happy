package api

// TODO: replace this with importing the happy-api repo

type App struct {
	AppName string `json:"app_name" validate:"required" gorm:"index:,unique,composite:metadata"`
}

var (
	EnvironmentMapping = map[string]uint8{
		"prod":    1,
		"staging": 2,
		"dev":     3,
		"rdev":    3,
	}
)

type AppMetadata struct {
	App
	Environment string `json:"environment"     validate:"required,valid_env" gorm:"index:,unique,composite:metadata"`
	Stack       string `json:"stack,omitempty"                               gorm:"default:'';not null;index:,unique,composite:metadata"`
} // @Name payload.AppMetadata

func NewAppMetadata(appName, env, stack string) *AppMetadata {
	return &AppMetadata{
		App:         App{AppName: appName},
		Environment: env,
		Stack:       stack,
	}
}

type ConfigKey struct {
	Key string `json:"key" validate:"required" gorm:"index:,unique,composite:metadata"`
}

type ConfigValue struct {
	ConfigKey
	Value string `json:"value" validate:"required"`
}

type AppConfigPayload struct {
	AppMetadata
	ConfigValue
}
type AppConfigLookupPayload struct {
	AppMetadata
	ConfigKey
} // @Name payload.AppConfigLookup

type CopyAppConfigPayload struct {
	App
	SrcEnvironment string `json:"source_environment"      validate:"required,valid_env" gorm:"index:,unique,composite:metadata"`
	SrcStack       string `json:"source_stack,omitempty"                                gorm:"default:'';not null;index:,unique,composite:metadata"`
	DstEnvironment string `json:"destination_environment" validate:"required,valid_env_dest" gorm:"index:,unique,composite:metadata"`
	DstStack       string `json:"destination_stack,omitempty"                           gorm:"default:'';not null;index:,unique,composite:metadata"`
	ConfigKey
} // @name payload.CopyAppConfig

type AppConfigDiffPayload struct {
	App
	SrcEnvironment string `json:"source_environment"      validate:"required,valid_env" gorm:"index:,unique,composite:metadata"`
	SrcStack       string `json:"source_stack,omitempty"                                gorm:"default:'';not null;index:,unique,composite:metadata"`
	DstEnvironment string `json:"destination_environment" validate:"required,valid_env_dest" gorm:"index:,unique,composite:metadata"`
	DstStack       string `json:"destination_stack,omitempty"                           gorm:"default:'';not null;index:,unique,composite:metadata"`
}
type AppConfig struct {
	AppConfigPayload
}
type WrappedAppConfigsWithCount struct {
	Records []AppConfig `json:"records"`
	Count   int         `json:"count" example:"1"`
}

type WrappedResolvedAppConfigsWithCount struct {
	Records []ResolvedAppConfig `json:"records"`
	Count   int                 `json:"count" example:"1"`
}

type ResolvedAppConfig struct {
	AppConfig
	Source string `json:"source"`
}

type WrappedResolvedAppConfig struct {
	Record *ResolvedAppConfig `json:"record"`
}

type WrappedAppConfig struct {
	Record AppConfig `json:"record"`
}

type ValidationError struct {
	FailedField string `json:"failed_field"` // the field that failed to be validated
	Tag         string `json:"tag" swaggerignore:"true"`
	Value       string `json:"value" swaggerignore:"true"`
	Type        string `json:"type" swaggerignore:"true"`
	Message     string `json:"message"` // a description of the error that occured
}

func NewAppConfigPayload(appName, env, stack, key, value string) *AppConfigPayload {
	return &AppConfigPayload{
		AppMetadata: *NewAppMetadata(appName, env, stack),
		ConfigValue: ConfigValue{
			Value: value,
			ConfigKey: ConfigKey{
				Key: key,
			},
		},
	}
}

func NewAppConfigLookupPayload(appName, env, stack, key string) *AppConfigLookupPayload {
	return &AppConfigLookupPayload{
		AppMetadata: *NewAppMetadata(appName, env, stack),
		ConfigKey: ConfigKey{
			Key: key,
		},
	}
}

func NewCopyAppConfigPayload(appName, srcEnv, srcStack, destEnv, destStack, key string) *CopyAppConfigPayload {
	return &CopyAppConfigPayload{
		App:            App{AppName: appName},
		SrcEnvironment: srcEnv,
		SrcStack:       srcStack,
		DstEnvironment: destEnv,
		DstStack:       destStack,
		ConfigKey: ConfigKey{
			Key: key,
		},
	}
}

func NewAppConfigDiffPayload(appName, srcEnv, srcStack, destEnv, destStack string) *AppConfigDiffPayload {
	return &AppConfigDiffPayload{
		App:            App{AppName: appName},
		SrcEnvironment: srcEnv,
		SrcStack:       srcStack,
		DstEnvironment: destEnv,
		DstStack:       destStack,
	}
}

type ConfigDiffResponse struct {
	MissingKeys []string `json:"missing_keys" example:"SOME_KEY,ANOTHER_KEY"`
}

type WrappedAppStacksWithCount struct {
	Records []*AppStack `json:"records"`
	Count   int         `json:"count" example:"1"`
}

type AppStack struct {
	AppStackPayload
}

type Enabler struct {
	Enabled *bool `json:"enabled" gorm:"default:true;not null;index:"`
}

type AppStackPayload struct {
	AppMetadata
	Enabler
}

type WrappedAppStack struct {
	Record *AppStack `json:"record"`
}
