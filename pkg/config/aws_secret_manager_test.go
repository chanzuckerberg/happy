package config

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetSecretValue(t *testing.T) {
	r := require.New(t)
	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(ctrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	},
		nil)

	awsSecretMgr := GetAwsSecretMgrWithClient(mock)
	secrets, err := awsSecretMgr.GetSecrets("test_arn")
	r.NoError(err)

	expected := &AwsSecretMgrSecrets{
		ClusterArn: "test_arn",
		Services: map[string]*RegistryConfig{"ecr_1": {
			Url: "test_url_1",
		}},
		Tfe: &TfeSecrets{
			Org: "tfe_org",
			Url: "tfe_url",
		},
	}
	r.Equal(expected, secrets)
}
