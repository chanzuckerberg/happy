package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
)

type StackECSCredentialsProvider struct{}

func (c StackECSCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	val := ctx.Value(request.AWSCredentialsContextKey{}).(request.AWSCredentials)
	return aws.Credentials{
		AccessKeyID:     val.AccessKeyID,
		SecretAccessKey: val.SecretAccessKey,
		SessionToken:    val.SessionToken,
	}, nil
}

func MakeCredentialProvider(ctx context.Context) aws.CredentialsProvider {
	return StackECSCredentialsProvider{}
}

type StackBackendECS struct{}

func getClient(ctx context.Context, payload model.AppStackPayload2) *ssm.Client {
	return ssm.New(ssm.Options{
		Region:      payload.AwsRegion,
		Credentials: MakeCredentialProvider(ctx),
	})
}

func (s *StackBackendECS) GetAppStacks(ctx context.Context, payload model.AppStackPayload2) ([]*model.AppStack, error) {
	client := getClient(ctx, payload)
	result, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/happy/%s/%s/stacklist", payload.AppName, payload.Environment)),
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not get parameter")
	}

	return convertParamToStacklist(*result.Parameter.Value, payload)
}
