package model

import (
	"fmt"
)

type App struct {
	AppName string `json:"app_name" query:"app_name" validate:"required" gorm:"index:,unique,composite:metadata" example:"testapp"`
}

var (
	EnvironmentMapping = map[string]uint8{
		"prod":    1,
		"staging": 2,
		"stage":   2,
		"sandbox": 2,
		"dev":     3,
		"rdev":    3,
	}
)

// NOTE: other happy related context will need to go here
// like happy config.json values for instance
type HappyContext struct {
	AWSContext
}

type AppMetadata struct {
	App
	HappyContext
	Environment string `json:"environment" query:"environment" validate:"required,valid_env" gorm:"index:,unique,composite:metadata" example:"rdev"`

	// in order to make this ON CONFLICT work we must not allow nulls for stack values
	// thus the stack column defaults to empty string and enforces NOT NULL
	Stack string `json:"stack,omitempty" query:"stack" gorm:"default:'';not null;index:,unique,composite:metadata" example:"my-stack"`
} // @Name payload.AppMetadata

func (a AppMetadata) String() string {
	return fmt.Sprintf("%s/%s/%s", a.App.AppName, a.Environment, a.Stack)
}

func NewAppMetadata(appName, env, stack string) *AppMetadata {
	return &AppMetadata{
		App:         App{AppName: appName},
		Environment: env,
		Stack:       stack,
	}
}
