package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var testFilePath = "testdata/test_config.yaml"

func TestNewHappConfig(t *testing.T) {
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

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)

			config, err := NewTestHappyConfig(t, testFilePath, testCase.env)
			r.NoError(err)

			r.Equal(config.TerraformVersion(), "0.13.5")
			r.Equal(config.DefaultEnv(), "rdev")
			r.Equal(config.App(), "test-app")
			r.Equal(config.SliceDefaultTag(), "branch-trunk")

			tasks, _ := config.GetTasks("migrate")
			r.Equal(tasks[0], "migrate_db_task_definition_arn")
			tasks, _ = config.GetTasks("delete")
			r.Equal(tasks[0], "delete_db_task_definition_arn")

			val := config.AwsProfile()
			r.Equal(testCase.wantAwsProfile, val)
			val = config.TerraformDirectory()
			r.Equal(testCase.wantTfDir, val)
			val = config.GetSecretArn()
			r.Equal(testCase.wantSecretArn, val)
			val = config.TaskLaunchType().String()
			r.Equal(testCase.wantTaskLaunchType, val)
		})
	}
}
