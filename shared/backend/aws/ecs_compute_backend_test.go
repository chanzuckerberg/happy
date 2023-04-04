package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../../config/testdata/test_config.yaml"
const testDockerComposePath = "../../config/testdata/docker-compose.yml"

func TestEcsComputeBackend(t *testing.T) {
	r := require.New(t)

	ctx := context.WithValue(context.Background(), util.CmdStartContextKey, time.Now())

	ctrl := gomock.NewController(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	secretsApi := interfaces.NewMockSecretsManagerAPI(ctrl)
	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secretsApi.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: &testVal,
			ARN:          aws.String("arn:aws:secretsmanager:region:accountid:secret:happy/env-happy-config-AB1234"),
		}, nil).AnyTimes()

	stsApi := interfaces.NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).
		Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil).AnyTimes()

	b, err := NewAWSBackend(ctx, happyConfig,
		WithAWSAccountID("1234567890"),
		WithSTSClient(stsApi),
		WithSecretsClient(secretsApi),
	)
	r.NoError(err)

	secret, secretArn, err := b.ComputeBackend.GetIntegrationSecret(ctx)
	r.NoError(err)

	r.IsType(&ECSComputeBackend{}, b.ComputeBackend)

	r.NotNil(secret)
	r.NotNil(secretArn)
	r.NotEmpty(*secretArn)
}
