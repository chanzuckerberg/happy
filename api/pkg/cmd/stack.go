package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
)

type StackManager interface {
	GetAppStacks(context.Context, model.AppStackPayload) ([]*model.AppStackResponse, error)
	// CreateOrUpdateAppStack(model.AppStackPayload) (*model.AppStack, error)
	// DeleteAppStack(model.AppStackPayload) (*model.AppStack, error)
}

type Stack struct{}

func MakeStack(db *store.DB) StackManager {
	return &Stack{}
}

func (s Stack) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
	// we cancel the context to close up any spun up goroutines
	// for this threads workspace_repo
	// TODO: we should probably cache these credentials/clients to TFE
	ctx, done := context.WithCancel(ctx)
	defer done()
	happyClient, err := request.MakeHappyClient(ctx, payload.AppName, payload.MakeEnvironmentContext(payload.Environment))
	if err != nil {
		return nil, errors.Wrap(err, "making happy client")
	}

	return happyClient.StackService.CollectStackInfo(ctx, payload.AppName, payload.ListAll)
}
