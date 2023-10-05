package model

type ConfigKey struct {
	Key string `json:"key" validate:"required" gorm:"index:,unique,composite:metadata" example:"SOME_KEY"`
}

type ConfigValue struct {
	ConfigKey
	Value string `json:"value" validate:"required" example:"some-value"`
}

type AppConfigPayload struct {
	AppMetadata
	ConfigValue
} // @Name payload.AppConfig

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
	SrcEnvironment string `json:"source_environment"          query:"source_environment"      validate:"required,valid_env"      gorm:"index:,unique,composite:metadata"`
	SrcStack       string `json:"source_stack,omitempty"      query:"source_stack"                                               gorm:"default:'';not null;index:,unique,composite:metadata"`
	DstEnvironment string `json:"destination_environment"     query:"destination_environment" validate:"required,valid_env_dest" gorm:"index:,unique,composite:metadata"`
	DstStack       string `json:"destination_stack,omitempty" query:"destination_stack"                                          gorm:"default:'';not null;index:,unique,composite:metadata"`
} // @name payload.AppConfigDiff

// @Description Object denoting which app config keys are missing from the destination env/stack
type ConfigDiffResponse struct {
	MissingKeys []string `json:"missing_keys" example:"SOME_KEY,ANOTHER_KEY"`
} // @Name response.ConfigDiff

// @Description App config key/value pair with additional metadata
type AppConfig struct {
	CommonDBFields
	AppConfigPayload
} // @Name response.AppConfig

// @Description App config key/value pair with additional metadata and "source"
type ResolvedAppConfig struct {
	AppConfig
	Source string `json:"source" example:"stack"`
} // @Name response.ResolvedAppConfig

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

func NewResolvedAppConfig(appConfig *AppConfig) *ResolvedAppConfig {
	source := "stack"
	if appConfig.Stack == "" {
		source = "environment"
	}
	return &ResolvedAppConfig{
		AppConfig: *appConfig,
		Source:    source,
	}
}

// @Description App config key/value pair wrapped in "record" key
type WrappedResolvedAppConfig struct {
	Record *ResolvedAppConfig `json:"record"`
} // @Name response.WrappedResolvedAppConfig

// @Description App config key/value pair wrapped in "record" key
type WrappedAppConfig struct {
	Record *AppConfig `json:"record"`
} // @Name response.WrappedAppConfig

type WrappedAppConfigsWithCount struct {
	Records []*AppConfig `json:"records"`
	Count   int          `json:"count" example:"1"`
} // @Name response.WrappedAppConfigsWithCount

type WrappedResolvedAppConfigsWithCount struct {
	Records []*ResolvedAppConfig `json:"records"`
	Count   int                  `json:"count" example:"1"`
} // @Name response.WrappedResolvedAppConfigsWithCount
