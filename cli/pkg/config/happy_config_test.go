package config

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

var testFilePath = "testdata/test_config.yaml"
var invalidTestFilePath = "testdata/test_config_invalid.yaml"

func TestNewHappyConfig(t *testing.T) {
	testData := []struct {
		env                   string
		wantAwsProfile        *string
		wantSecretId          string
		wantTfDir             string
		wantTaskLaunchType    string
		wantAutorunMigrations bool
	}{
		{"rdev", aws.String("test-dev"), "happy/env-rdev-config", ".happy/terraform/envs/rdev", "EC2", true},
		{"stage", aws.String("test-stage"), "happy/env-stage-config", ".happy/terraform/envs/stage", "FARGATE", false},
		{"prod", aws.String("test-prod"), "happy/env-prod-config", ".happy/terraform/envs/prod", "FARGATE", false},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)

			config, err := NewTestHappyConfig(t, testFilePath, testCase.env)
			r.NoError(err)

			r.Equal(config.DefaultEnv(), "rdev")
			r.Equal(config.App(), "test-app")
			r.Equal(config.SliceDefaultTag(), "branch-trunk")
			r.Equal(testCase.wantAutorunMigrations, config.AutoRunMigrations())

			tasks, _ := config.GetTasks("migrate")
			r.Equal(tasks[0], "migrate_db_task_definition_arn")
			tasks, _ = config.GetTasks("delete")
			r.Equal(tasks[0], "delete_db_task_definition_arn")

			awsProfval := config.AwsProfile()
			r.Equal(testCase.wantAwsProfile, awsProfval)
			val := config.TerraformDirectory()
			r.Equal(testCase.wantTfDir, val)
			val = config.GetSecretId()
			r.Equal(testCase.wantSecretId, val)
			val = config.TaskLaunchType().String()
			r.Equal(testCase.wantTaskLaunchType, val)
		})
	}
}

func TestProfile(t *testing.T) {
	r := require.New(t)

	var nilProfile *Profile
	r.Equal("*", nilProfile.Get())

	otherProfile := Profile("foobarother")
	r.Equal("foobarother", otherProfile.Get())
}

func TestMissingDefaultEnvConfig(t *testing.T) {
	r := require.New(t)
	_, err := NewTestHappyConfig(t, invalidTestFilePath, "")
	r.Error(err)
}

func TestDefaultEnvPriority(t *testing.T) {
	r := require.New(t)

	config, err := NewTestHappyConfig(t, testFilePath, "")
	r.NoError(err)
	r.Equal(config.GetEnv(), config.DefaultEnv())

	testEnv := "rdev"
	config, err = NewTestHappyConfig(t, testFilePath, testEnv)
	r.NoError(err)
	r.Equal(config.GetEnv(), testEnv)
}
