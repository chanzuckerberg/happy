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
<<<<<<< HEAD
	// we cancel the context to close up any spun up goroutines
	// for this threads workspace_repo
	// TODO: we should probably cache these credentials/clients to TFE
	ctx, done := context.WithCancel(ctx)
	defer done()
	happyClient, err := request.MakeHappyClient(ctx, payload.AppName, payload.MakeEnvironmentContext(payload.Environment))
=======
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
	ctx, done := context.WithCancel(ctx)
	defer done()
	workspace_repo.StartTFCWorkerPool(ctx)

	wg := sync.WaitGroup{}

	stackInfos := make([]*model.AppStackResponse, len(stacklist))
	for i, stackName := range stacklist {
		wg.Add(1)
		go func(id int, stackName string) {
			defer wg.Done()

			stackInfo, err := getStackInfo(ctx, payload, stackName, workspaceRepo)
			if err != nil {
				stackInfo.Error = err.Error()
			}
			stackInfos[id] = stackInfo
		}(i, stackName)
	}
	wg.Wait()

	return stackInfos, nil
}

func getStackInfo(ctx context.Context, payload model.AppStackPayload, stackName string, workspaceRepo workspace_repo.WorkspaceRepoIface) (*model.AppStackResponse, error) {
	stack := &model.AppStackResponse{
		AppMetadata: *model.NewAppMetadata(payload.AppName, payload.Environment, stackName),
	}
	workspace, err := workspaceRepo.GetWorkspace(ctx, fmt.Sprintf("%s-%s", payload.AppMetadata.Environment, stackName))
>>>>>>> origin/main
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
