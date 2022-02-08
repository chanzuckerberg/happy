package workspace_repo

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"

func TestNewWorkspaceRepoErrorNoTFEToken(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(ctrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	}, nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)

	bootstrap := &config.Bootstrap{HappyConfigPath: testFilePath}
	happyConfig, err := config.NewHappyConfigWithSecretsBackend(bootstrap, awsSecretMgr)
	r.NoError(err)

	_, err = NewWorkspaceRepo(happyConfig, "foo", "bar")
	r.True(err == nil || strings.Contains(err.Error(), "please set env var TFE_TOKEN"))
}
