package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	compute_backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
)

type StackBackendECS struct{}

func getSSMClient(ctx context.Context, payload model.AppStackPayload) *ssm.Client {
	return ssm.New(ssm.Options{
		Region:      payload.AwsRegion,
		Credentials: request.MakeCredentialProvider(ctx),
	})
}

func getSecretsManagerClient(ctx context.Context, payload model.AppStackPayload) *secretsmanager.Client {
	return secretsmanager.New(secretsmanager.Options{
		Region:      payload.AwsRegion,
		Credentials: request.MakeCredentialProvider(ctx),
	})
}

func (s *StackBackendECS) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
	ssmClient := getSSMClient(ctx, payload)
	result, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/happy/%s/%s/stacklist", payload.AppName, payload.Environment)),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ParameterNotFound") {
			return []*model.AppStackResponse{}, nil
		}
		return nil, errors.Wrap(err, "could not get parameter")
	}

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
	integrationSecret, _, err := computeBackend.GetIntegrationSecret(ctx)
	if err != nil {
		return nil, err
	}

	stacklist, err := parseParamToStacklist(*result.Parameter.Value)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json")
	}

	return enrichStacklistMetadata(ctx, stacklist, payload, integrationSecret)
}
