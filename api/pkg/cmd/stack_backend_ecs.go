package cmd

import (
	"context"
	"fmt"

	compute_backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
)

type StackBackendECS struct{}

func (s *StackBackendECS) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
	envCtx := config.EnvironmentContext{
		EnvironmentName: payload.Environment,
		AWSProfile:      &payload.AwsProfile,
		AWSRegion:       &payload.AwsRegion,
		SecretId:        payload.SecretId,
		TaskLaunchType:  util.LaunchType(payload.TaskLaunchType),
	}
	b, err := compute_backend.NewAWSBackend(ctx, envCtx)
	if err != nil {
		return nil, errors.Wrap(err, "could not create backend")
	}

	computeBackend := compute_backend.ECSComputeBackend{
		Backend:  b,
		SecretId: payload.SecretId,
	}

	paramOutput, err := computeBackend.GetParam(ctx, fmt.Sprintf("/happy/%s/%s/stacklist", payload.AppName, payload.Environment))
	if err != nil {
		return nil, err
	}

	integrationSecret, _, err := computeBackend.GetIntegrationSecret(ctx)
	if err != nil {
		return nil, err
	}

	stacklist, err := parseParamToStacklist(paramOutput)
	if err != nil {
		return nil, err
	}

	return enrichStacklistMetadata(ctx, stacklist, payload, integrationSecret)
}
