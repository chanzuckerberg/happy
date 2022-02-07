package config

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var testFilePath = "testdata/test_config.yaml"

func TestNewHappyConfig(t *testing.T) {
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

	testData := []struct {
		env                string
		wantAwsProfile     string
		wantSecretArn      string
		wantTfDir          string
		wantTaskLaunchType string
	}{
		{"rdev", "test-dev", "happy/env-rdev-config", ".happy/terraform/envs/rdev", "FARGATE"},
		{"stage", "test-stage", "happy/env-stage-config", ".happy/terraform/envs/stage", "FARGATE"},
		{"prod", "test-prod", "happy/env-prod-config", ".happy/terraform/envs/prod", "FARGATE"},
	}
	secretValue, err := awsSecretMgr.GetSecrets("cluster_arn")
	r.NoError(err)
	r.Equal(secretValue.GetClusterArn(), "test_arn")

	for _, testCase := range testData {
		config, err := NewTestHappyConfig(t, testFilePath, testCase.env, awsSecretMgr)
		r.NoError(err)

		r.Equal(config.TerraformVersion(), "0.13.5")
		r.Equal(config.GetEnv(), testCase.env)
		r.Equal(config.App(), "test-app")
		r.Equal(config.SliceDefaultTag(), "branch-trunk")

		slices, _ := config.GetSlices()
		r.Equal(slices["backend"].BuildImages[0], "backend")
		r.Equal(slices["frontend"].BuildImages[0], "frontend")

		tasks, _ := config.GetTasks("migrate")
		r.Equal(tasks[0], "migrate_db_task_definition_arn")
		tasks, _ = config.GetTasks("delete")
		r.Equal(tasks[0], "delete_db_task_definition_arn")

		val := config.AwsProfile()
		r.Equal(val, testCase.wantAwsProfile)
		val = config.TerraformDirectory()
		r.Equal(val, testCase.wantTfDir)
		val = config.GetSecretArn()
		r.Equal(val, testCase.wantSecretArn)
		val = config.TaskLaunchType()
		r.Equal(val, testCase.wantTaskLaunchType)

		serviceRegistries := config.GetRdevServiceRegistries()
		r.True(len(serviceRegistries) > 0)
	}
}
