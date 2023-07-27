package request

import (
	"context"

	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
)

type HappyClient struct {
	StackService *stack.StackService
	AWSBackend   *backend.Backend
}

func MakeHappyClient(ctx context.Context, appName string, envCtx config.EnvironmentContext) (*HappyClient, error) {
	awsBackend, err := backend.NewAWSBackend(ctx, envCtx,
		backend.WithNewAWSConfigOption(configv2.WithRegion(*envCtx.AWSRegion)),
		backend.WithNewAWSConfigOption(configv2.WithCredentialsProvider(MakeCredentialProvider(ctx))),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct an AWS backend")
	}

	workspaceRepo := workspace_repo.NewWorkspaceRepo(awsBackend.Conf().GetTfeUrl(), awsBackend.Conf().GetTfeOrg()).
		WithTFEToken(setup.GetConfiguration().TFE.Token)
	workspace_repo.StartTFCWorkerPool(ctx)

	stackService := stack.NewStackService().
		WithApp(envCtx.EnvironmentName, appName).
		WithBackend(awsBackend).
		WithWorkspaceRepo(workspaceRepo)

	return &HappyClient{
		StackService: stackService,
		AWSBackend:   awsBackend,
	}, nil
}
