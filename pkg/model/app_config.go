package model

import "gorm.io/gorm"

type ConfigKey struct {
	Key string `json:"key"   validate:"required" gorm:"index:,unique,composite:metadata"`
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
}

type CopyAppConfigPayload struct {
	AppName        string `json:"app_name"                validate:"required" gorm:"index:,unique,composite:metadata"`
	SrcEnvironment string `json:"source_environment"      validate:"required" gorm:"index:,unique,composite:metadata"`
	SrcStack       string `json:"source_stack,omitempty"                      gorm:"default:'';not null;index:,unique,composite:metadata"`
	DstEnvironment string `json:"destination_environment" validate:"required" gorm:"index:,unique,composite:metadata"`
	DstStack       string `json:"destination_stack,omitempty"                 gorm:"default:'';not null;index:,unique,composite:metadata"`
	ConfigKey
}

type AppConfigResponse struct {
	AppConfig
	Source string `json:"source"`
}

type AppConfig struct {
	gorm.Model
	AppConfigPayload
}

func NewAppConfigPayload(appName, env, stack, key, value string) *AppConfigPayload {
	return &AppConfigPayload{
		AppMetadata: AppMetadata{
			AppName:     appName,
			Environment: env,
			Stack:       stack,
		},
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
		AppMetadata: AppMetadata{
			AppName:     appName,
			Environment: env,
			Stack:       stack,
		},
		ConfigKey: ConfigKey{
			Key: key,
		},
	}
}

func NewCopyAppConfigPayload(appName, srcEnv, srcStack, destEnv, destStack, key string) *CopyAppConfigPayload {
	return &CopyAppConfigPayload{
		AppName:        appName,
		SrcEnvironment: srcEnv,
		SrcStack:       srcStack,
		DstEnvironment: destEnv,
		DstStack:       destStack,
		ConfigKey: ConfigKey{
			Key: key,
		},
	}
}
