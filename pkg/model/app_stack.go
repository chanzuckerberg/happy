package model

import "gorm.io/gorm"

type Enabler struct {
	Enabled *bool `json:"enabled" gorm:"default:true;not null;index:"`
}

type AppStack struct {
	gorm.Model `swaggerignore:"true"`
	AppStackPayload
}

type AppStackPayload struct {
	AppMetadata
	Enabler
} // @Name payload.AppStackPayload

func MakeAppStack(appName, env, stack string, enabled bool) AppStack {
	return AppStack{
		AppStackPayload: MakeAppStackPayload(appName, env, stack, enabled),
	}
}

func MakeAppStackPayload(appName, env, stack string, enabled bool) AppStackPayload {
	return AppStackPayload{
		AppMetadata: *NewAppMetadata(appName, env, stack),
		Enabler:     Enabler{&enabled},
	}
}
