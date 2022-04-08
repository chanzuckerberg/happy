package config

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

var testFilePath = "testdata/test_config.yaml"

func TestNewHappConfig(t *testing.T) {
	testData := []struct {
		env                   string
		wantAwsProfile        *string
		wantSecretArn         string
		wantTfDir             string
		wantTaskLaunchType    string
		wantAutorunMigrations bool
	}{
		{"rdev", aws.String("test-dev"), "happy/env-rdev-config", ".happy/terraform/envs/rdev", "EC2", true},
		{"stage", aws.String("test-stage"), "happy/env-stage-config", ".happy/terraform/envs/stage", "FARGATE", false},
		{"prod", aws.String("test-prod"), "happy/env-prod-config", ".happy/terraform/envs/prod", "FARGATE", false},
	}

	r := require.New(t)

	targetPlatforms := map[string]bool{}
	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			config, err := NewTestHappyConfig(t, testFilePath, testCase.env)
			r.NoError(err)

			r.Equal(config.TerraformVersion(), "0.13.5")
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
			val = config.GetSecretArn()
			r.Equal(testCase.wantSecretArn, val)
			val = config.TaskLaunchType().String()
			r.Equal(testCase.wantTaskLaunchType, val)
			targetPlatforms[config.GetTargetContainerPlatform()] = true
		})
	}

	r.Len(targetPlatforms, 2)
}

func TestProfile(t *testing.T) {
	r := require.New(t)

	var nilProfile *Profile
	r.Equal("*", nilProfile.Get())

	otherProfile := Profile("foobarother")
	r.Equal("foobarother", otherProfile.Get())
}
