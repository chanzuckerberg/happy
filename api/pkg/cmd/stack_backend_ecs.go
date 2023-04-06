package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
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

	secretsClient := getSecretsManagerClient(ctx, payload)
	out, err := secretsClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &payload.SecretId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not get integration secret at %s", payload.SecretId)
	}

	secret := &config.IntegrationSecret{}
	err = json.Unmarshal([]byte(*out.SecretString), secret)
	if err != nil {
		return nil, errors.Wrap(err, "could not json parse integraiton secret")
	}

	stacklist, err := parseParamToStacklist(*result.Parameter.Value)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json")
	}

	return enrichStacklistMetadata(ctx, stacklist, payload, secret)
}
