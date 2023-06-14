package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
)

type StackManager interface {
	GetAppStacks(context.Context, model.AppStackPayload) ([]*model.AppStackResponse, error)
	// CreateOrUpdateAppStack(model.AppStackPayload) (*model.AppStack, error)
	// DeleteAppStack(model.AppStackPayload) (*model.AppStack, error)
}

type Stack struct {
	// TODO: see note below about eventually read and writing stack information to a database
	// keeping this as a reminder
	//db StackManager
}

func MakeStack(db *dbutil.DB) StackManager {
	return &Stack{
		// DB is not currently used since this is currently just a read interface for the old data locations
		// but we should keep this here so it's easy to set up later when we want to move the data
		//db: MakeStackBackendDB(db),
	}
}

func (s Stack) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
	happyClient, err := request.MakeHappyClient(ctx, payload.AppName, payload.MakeEnvironmentContext(payload.Environment))
	if err != nil {
		return nil, errors.Wrap(err, "making happy client")
	}

	stacks, err := happyClient.StackService.CollectStackInfo(ctx, false, payload.AppName)
	if err != nil {
		return nil, errors.Wrapf(err, "collecting stack info")
	}

	resp := make([]*model.AppStackResponse, len(stacks))
	for i, stack := range stacks {
		resp[i] = &model.AppStackResponse{
			AppMetadata:   *model.NewAppMetadata(stack.App, payload.Environment, stack.Name),
			StackMetadata: stack,
		}
	}
	return resp, nil
}
