package orchestrator

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"
const testDockerComposePath = "../config/testdata/docker-compose.yml"

func TestNewOrchestrator(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	secrets := mocks.NewMockSecretsManagerAPI(ctrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secrets.EXPECT().GetSecretValueWithContext(ctx, gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretBinary: []byte(testVal),
	}, nil)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
	r.NoError(err)

	backend, err := backend.NewAWSBackend(ctx, happyConfig, backend.WithSecretsClient(secrets))
	r.NoError(err)

	orchestrator := NewOrchestrator(backend)
	r.NotNil(orchestrator)
	err = orchestrator.Shell("frontend", "")
	r.Error(err)
}
