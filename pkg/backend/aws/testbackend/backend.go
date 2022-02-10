package testbackend

import (
	context "context"

	"github.com/aws/aws-sdk-go/aws"
	secretsmanager "github.com/aws/aws-sdk-go/service/secretsmanager"
	sts "github.com/aws/aws-sdk-go/service/sts"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
)

// NewBackend will return a test backend with
// some sane default mocks
// These can be overriden as needed by callers
func NewBackend(
	ctx context.Context,
	ctrl *gomock.Controller,
	conf *config.HappyConfig,
	opts ...backend.AWSBackendOption) (*backend.Backend, error) {

	// first set our own defaults
	secrets := NewMockSecretsManagerAPI(ctrl)
	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secrets.EXPECT().GetSecretValueWithContext(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: &testVal,
		}, nil)

	stsApi := NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentityWithContext(gomock.Any(), gomock.Any()).
		Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil)

	// then add provided
	// note how the user-provided ones are the last in the slice and therefore they will override our defaults
	combinedOpts := []backend.AWSBackendOption{
		backend.WithSecretsClient(secrets),
		backend.WithSTSClient(stsApi),
	}
	if opts != nil {
		combinedOpts = append(combinedOpts, opts...)
	}

	return backend.NewAWSBackend(ctx, conf, combinedOpts...)
}
