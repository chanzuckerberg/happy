package testbackend

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/golang/mock/gomock"
)

// NewBackend will return a test backend with
// some sane default mocks
// These can be overriden as needed by callers
func NewBackend(
	ctx context.Context,
	ctrl *gomock.Controller,
	environmentContext config.EnvironmentContext,
	opts ...backend.AWSBackendOption) (*backend.Backend, error) {
	// first set our own defaults
	secrets := interfaces.NewMockSecretsManagerAPI(ctrl)
	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secrets.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: &testVal,
			ARN:          aws.String("arn:aws:secretsmanager:region:accountid:secret:happy/env-happy-config-AB1234"),
		}, nil).AnyTimes()

	stsApi := interfaces.NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).
		Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil).AnyTimes()

	// by default, prevent all calls unless specifically overriden
	cwl := interfaces.NewMockGetLogEventsAPIClient(ctrl)
	ecs := interfaces.NewMockECSAPI(ctrl)
	ec2 := interfaces.NewMockEC2API(ctrl)
	ecr := interfaces.NewMockECRAPI(ctrl)

	// then add provided
	// note how the user-provided ones are the last in the slice, and therefore they will override our defaults
	combinedOpts := []backend.AWSBackendOption{
		backend.WithAWSAccountID("1234567890"),
		backend.WithSecretsClient(secrets),
		backend.WithSTSClient(stsApi),
		backend.WithGetLogEventsAPIClient(cwl),
		backend.WithECSClient(ecs),
		backend.WithEC2Client(ec2),
		backend.WithECRClient(ecr),
	}

	// apply opts
	if opts != nil {
		combinedOpts = append(combinedOpts, opts...)
	}

	return backend.NewAWSBackend(ctx, environmentContext, combinedOpts...)
}
