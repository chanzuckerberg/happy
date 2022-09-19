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
