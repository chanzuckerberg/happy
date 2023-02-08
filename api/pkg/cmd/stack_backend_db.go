package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type StackBackendDB struct {
	DB *dbutil.DB
}

func MakeStackBackendDB(db *dbutil.DB) *StackBackendDB {
	return &StackBackendDB{
		DB: db,
	}
}

func (s *StackBackendDB) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStack, error) {
	db := s.DB.GetDB()
	stack := &model.AppStack{AppMetadata: payload.AppMetadata}
	stacks := []*model.AppStack{}
	res := db.Where(stack).Find(&stacks)
	return stacks, errors.Wrapf(res.Error, "unable to get app stacks for %s", stack.AppMetadata)
}

func (s *StackBackendDB) CreateOrUpdateAppStack(payload model.AppStackPayload) (*model.AppStack, error) {
	db := s.DB.GetDB()
	stack := model.NewAppStackFromAppStackPayload(payload)
	res := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "app_name"},
			{Name: "environment"},
			{Name: "stack"},
		},
		UpdateAll: true,
	}).Create(&stack)

	return stack, errors.Wrapf(res.Error, "unable to create app stack %s", payload.AppMetadata)
}

func (s *StackBackendDB) DeleteAppStack(payload model.AppStackPayload) (*model.AppStack, error) {
	db := s.DB.GetDB()
	record := &model.AppStack{}
	res := db.Clauses(clause.Returning{}).
		Where("app_name = ? AND environment = ? AND stack = ?",
			payload.AppName,
			payload.Environment,
			payload.Stack,
		).Delete(record)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return record, nil
}
