package model

import (
	"github.com/chanzuckerberg/happy/shared/k8s"
	"gorm.io/gorm"
)

type AppStack struct {
	gorm.Model `swaggerignore:"true"`
	AppStackPayload
}

type AppStackPayload struct {
	AppMetadata
} // @Name payload.AppStackPayload

type WrappedAppStacksWithCount struct {
	Records []*AppStack `json:"records"`
	Count   int         `json:"count" example:"1"`
} // @Name response.WrappedAppStacksWithCount

type WrappedAppStack struct {
	Record *AppStack `json:"record"`
} // @Name response.WrappedAppStack

func MakeAppStack(appName, env, stack string) AppStack {
	return AppStack{
		AppStackPayload: MakeAppStackPayload(appName, env, stack),
	}
}

func MakeAppStackPayload(appName, env, stack string) AppStackPayload {
	return AppStackPayload{
		AppMetadata: *NewAppMetadata(appName, env, stack),
	}
}

type AppStackPayload2 struct {
	AppName        string         `json:"app_name" validate:"required" gorm:"index:,unique,composite:metadata"`
	Environment    string         `json:"environment" validate:"required,valid_env" gorm:"index:,unique,composite:metadata"`
	AwsProfile     string         `json:"aws_profile" validate:"required"`
	AwsRegion      string         `json:"aws_region" validate:"required"`
	TaskLaunchType string         `json:"task_launch_type" validate:"required"`
	K8SConfig      *k8s.K8SConfig `json:"k8s"`
}
