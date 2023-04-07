package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StackManager interface {
	GetAppStacks(context.Context, model.AppStackPayload) ([]*model.AppStackResponse, error)
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

func (s Stack) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
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

func parseParamToStacklist(paramOutput string) ([]string, error) {
	var stacklist []string
	err := json.Unmarshal([]byte(paramOutput), &stacklist)
	return stacklist, errors.Wrap(err, "could not parse json")
}

func enrichStacklistMetadata(ctx context.Context, stacklist []string, payload model.AppStackPayload, integrationSecret *config.IntegrationSecret) ([]*model.AppStackResponse, error) {
	workspaceRepo := workspace_repo.NewWorkspaceRepo(
		integrationSecret.Tfe.Url,
		integrationSecret.Tfe.Org,
	).WithTFEToken(setup.GetConfiguration().TFE.Token)
	wg := sync.WaitGroup{}

	stacks := []*model.AppStackResponse{}

	for _, stackName := range stacklist {
		wg.Add(1)
		go func(stackName string) {
			defer wg.Done()

			stack := &model.AppStackResponse{
				AppMetadata: *model.NewAppMetadata(payload.AppName, payload.Environment, stackName),
			}
			// the error handling below is not the typical "if err != nil { return err }"  because
			// if this errors we still want to return the stack, it just won't have all the fields populated
			workspace, err := workspaceRepo.GetWorkspace(ctx, fmt.Sprintf("%s-%s", payload.AppMetadata.Environment, stackName))
			if err == nil {
				stack.WorkspaceUrl = workspace.GetWorkspaceUrl()
				stack.WorkspaceStatus = workspace.GetCurrentRunStatus(ctx)
				stack.Endpoints = map[string]string{}

				// the error handling below is not the typical "if err != nil { return err }"  because
				// if this errors we still want to return the stack, it just won't have the Endpoints field populated
				endpoints, err := workspace.GetEndpoints(ctx)
				if err == nil {
					stack.Endpoints = endpoints
				}
			}
			stacks = append(stacks, stack)
		}(stackName)
	}
	wg.Wait()

	return stacks, nil
}
