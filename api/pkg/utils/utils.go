package utils

import (
	"context"
	"encoding/json"
	"os"

	"github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
)

func StructToMap(payload interface{}) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(payload)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}
	return inInterface, nil
}

type HappyClient struct {
	StackService *stack.StackService
	AWSBackend   *aws.Backend
}

func MakeHappyClient(ctx context.Context, a model.AppMetadata) (*HappyClient, error) {
	awsBackend, err := a.MakeAWSBackend(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "making AWS backend")
	}

	// TODO: we need to enforce the TFE_TOKEN exists on boot
	// TODO: do we want to pass in the TFE_TOKEN instead of building with a TFE_TOKEN as part of the environment?
	tfeToken := os.Getenv("TFE_TOKEN")
	if tfeToken == "" {
		return nil, errors.New("TFE_TOKEN is required; this should never happen")
	}
	workspaceRepo := workspace_repo.NewWorkspaceRepo(awsBackend.Conf().GetTfeUrl(), awsBackend.Conf().GetTfeOrg()).
		WithTFEToken(tfeToken)

	stackService := stack.NewStackService(a.Environment, a.AppName).
		WithBackend(awsBackend).
		WithWorkspaceRepo(workspaceRepo)

	return &HappyClient{
		StackService: stackService,
		AWSBackend:   awsBackend,
	}, nil
}
