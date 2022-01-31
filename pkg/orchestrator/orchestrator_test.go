package orchestrator

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var testFilePath = "../config/test_config.yaml"

func TestNewOrchestrator(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(ctrl)

	testVal := "{\"cluster_arn\":\"test_arn\",\"ecrs\":{\"ecr_1\":{\"url\":\"test_url_1\"}}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	},
		nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)
	r.NotNil(awsSecretMgr)

	happyConfig, err := config.NewHappyConfig(testFilePath, "rdev")

	orchestrator := NewOrchestrator(happyConfig, nil)
	r.NotNil(orchestrator)
	err = orchestrator.Shell("frontend", "")
	r.Nil(err)
}
