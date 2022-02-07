package orchestrator

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"

func TestNewOrchestrator(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(ctrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	}, nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)
	r.NotNil(awsSecretMgr)

	happyConfig, err := NewTestHappyConfig(t, testFilePath, "rdev", awsSecretMgr)
	r.NoError(err)

	taskRunner := backend.GetAwsEcs(happyConfig)
	orchestrator := NewOrchestrator(happyConfig, taskRunner)
	r.NotNil(orchestrator)
	err = orchestrator.Shell("frontend", "")
	r.Error(err)
}

// generates a test happy config
// only use in tests
func NewTestHappyConfig(
	t *testing.T,
	testFilePath string,
	env string,
	awsSecretMgr config.SecretsBackend,
) (config.HappyConfig, error) {
	b := &config.Bootstrap{
		Env:             env,
		HappyConfigPath: testFilePath,
	}
	return config.NewHappyConfigWithSecretsBackend(b, awsSecretMgr)
}
