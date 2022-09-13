package cmd

import (
	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/pkg/errors"
)

type Stack interface {
	CreateAppStack(model.AppStackPayload) (*model.AppStack, error)
	GetAppStacks(model.AppStackPayload) ([]*model.AppStack, error)
	UpdateAppStack(model.AppStackPayload) (*model.AppStack, error)
}

type dbStack struct {
	DB *dbutil.DB
}

func MakeStack(db *dbutil.DB) Stack {
	return &dbStack{
		DB: db,
	}
}

func (s *dbStack) CreateAppStack(payload model.AppStackPayload) (*model.AppStack, error) {
	db := s.DB.GetDB()
	stack := &model.AppStack{AppStackPayload: payload}
	res := db.Create(stack)
	return stack, errors.Wrapf(res.Error, "unable to create app stack %s", payload.AppMetadata)
}

func (s *dbStack) GetAppStacks(payload model.AppStackPayload) ([]*model.AppStack, error) {
	db := s.DB.GetDB()
	stack := &model.AppStack{AppStackPayload: payload}
	stacks := []*model.AppStack{}
	if stack.Enabled == nil {
		defaultTrue := true
		stack.Enabled = &defaultTrue
	}
	res := db.Where(stack).Find(&stacks)
	return stacks, errors.Wrapf(res.Error, "unable to get app stacks for %s", stack.AppMetadata)
}

func (s *dbStack) UpdateAppStack(payload model.AppStackPayload) (*model.AppStack, error) {
	db := s.DB.GetDB()
	stack := &model.AppStack{AppStackPayload: payload}
	res := db.Model(&stack).
		Where("app_name = ? AND environment = ? AND stack = ? AND enabled = ?",
			payload.AppName,
			payload.Environment,
			payload.Stack,
			payload.Enabled).
		Updates(stack)
	if res.RowsAffected == 0 {
		return nil, errors.Errorf("found no rows to update for %s", stack.AppMetadata)
	}
	return stack, errors.Wrapf(res.Error, "unable to update app stack for %s", stack.AppMetadata)
}
