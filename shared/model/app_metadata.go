package model

import "fmt"

type App struct {
	AppName string `json:"app_name" query:"app_name" validate:"required" example:"testapp"`
}

var (
	EnvironmentMapping = map[string]uint8{
		"prod":    1,
		"staging": 2,
		"stage":   2,
		"sandbox": 2,
		"dev":     3,
		"rdev":    3,
		"ldev":    3,
	}
)

type AppMetadata struct {
	App
	Environment string `json:"environment" query:"environment" validate:"required,valid_env" example:"rdev"`

	// in order to make this ON CONFLICT work we must not allow nulls for stack values
	// thus the stack column defaults to empty string and enforces NOT NULL
	Stack string `json:"stack,omitempty" query:"stack" example:"my-stack"`
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
