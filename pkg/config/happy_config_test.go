package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testFilePath = "testdata/test_config.yaml"

func TestNewHappConfig(t *testing.T) {
	// Setup
	r := require.New(t)

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

	// Run tests
	for _, testCase := range testData {
		config, err := NewHappyConfig(testFilePath, testCase.env)
		r.NoError(err)

		r.Equal(config.TerraformVersion(), "0.13.5")
		r.Equal(config.DefaultEnv(), "rdev")
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
	}
}
