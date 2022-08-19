package model

type AppMetadata struct {
	AppName     string `json:"app_name"    validate:"required" gorm:"index:,unique,composite:metadata"`
	Environment string `json:"environment" validate:"required" gorm:"index:,unique,composite:metadata"`
	Stack       string `json:"stack,omitempty"                 gorm:"default:'';not null;index:,unique,composite:metadata"`
}
