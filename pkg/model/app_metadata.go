package model

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
