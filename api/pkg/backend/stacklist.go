package backend

import "github.com/chanzuckerberg/happy/shared/model"

type StacklistIface interface {
	GetAppStacks(model.AppStackPayload) ([]*model.AppStack, error)
	CreateOrUpdateAppStack(model.AppStackPayload) (*model.AppStack, error)
	DeleteAppStack(model.AppStackPayload) (*model.AppStack, error)
}
