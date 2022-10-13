package model

import "fmt"

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
	Environment string `json:"environment" validate:"required,valid_env" gorm:"index:,unique,composite:metadata"`

	// in order to make this ON CONFLICT work we must not allow nulls for stack values
	// thus the stack column defaults to empty string and enforces NOT NULL
	Stack string `json:"stack,omitempty" gorm:"default:'';not null;index:,unique,composite:metadata"`
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
