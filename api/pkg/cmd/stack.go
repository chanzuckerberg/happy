package cmd

import (
	"context"
	"encoding/json"

	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StackManager interface {
	GetAppStacks(context.Context, model.AppStackPayload) ([]*model.AppStack, error)
	// CreateOrUpdateAppStack(model.AppStackPayload) (*model.AppStack, error)
	// DeleteAppStack(model.AppStackPayload) (*model.AppStack, error)
}

type Stack struct {
	db  StackManager
	ecs StackManager
	eks StackManager
}

func MakeStack(db *dbutil.DB) StackManager {
	return &Stack{
		// DB is not currently used since this is currently just a read interface for the old data locations
		// but we should keep this here so it's easy to set up later when we want to move the data
		db:  MakeStackBackendDB(db),
		ecs: &StackBackendECS{},
		eks: &StackBackendEKS{},
	}
}

func (s Stack) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStack, error) {
	switch payload.TaskLaunchType {
	case "k8s":
		return s.eks.GetAppStacks(ctx, payload)
	case "fargate":
		return s.ecs.GetAppStacks(ctx, payload)
	default:
		logrus.Fatal("Must specify a Launch Type as either k8s or fargate")
	}
	return nil, nil
}

func convertParamToStacklist(paramOutput string, payload model.AppStackPayload) ([]*model.AppStack, error) {
	var stacklist []string
	err := json.Unmarshal([]byte(paramOutput), &stacklist)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json")
	}

	stacks := []*model.AppStack{}
	for _, stackName := range stacklist {
		appStack := model.MakeAppStack(payload.AppName, payload.Environment, stackName)
		stacks = append(stacks, &appStack)
	}

	return stacks, nil
}
