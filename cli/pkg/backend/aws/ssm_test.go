package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSSMParams(t *testing.T) {
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

	stsApi := interfaces.NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).
		Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil).AnyTimes()

	secretsApi := interfaces.NewMockSecretsManagerAPI(ctrl)
	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secretsApi.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: &testVal,
			ARN:          aws.String("arn:aws:secretsmanager:region:accountid:secret:happy/env-happy-config-AB1234"),
		}, nil).AnyTimes()

	testParamStoreData := "value"
	ssmApi := interfaces.NewMockSSMAPI(ctrl)
	ssmApi.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(&ssm.GetParameterOutput{Parameter: &ssmtypes.Parameter{Value: &testParamStoreData}}, nil)
	ssmApi.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(&ssm.PutParameterOutput{}, nil)
	b, err := NewAWSBackend(ctx, happyConfig,
		WithAWSAccountID("1234567890"),
		WithSTSClient(stsApi),
		WithSSMClient(ssmApi),
		WithSecretsClient(secretsApi),
	)
	r.NoError(err)

	param, err := b.ComputeBackend.GetParam(ctx, "/param")
	r.NoError(err)
	r.NotEmpty(param)

	err = b.ComputeBackend.WriteParam(ctx, "/param", "value")
	r.NoError(err)
}
