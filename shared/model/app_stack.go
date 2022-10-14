package model

import (
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
