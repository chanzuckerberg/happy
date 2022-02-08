package stack_mgr

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

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	r := require.New(t)
	ctrl := gomock.NewController(t)
	secrets := mocks.NewMockSecretsManagerAPI(ctrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secrets.EXPECT().GetSecretValueWithContext(gomock.Any(), gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	}, nil)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	config, err := config.NewHappyConfig(ctx, bootstrapConfig)
	r.NoError(err)

	dataMap := map[string]string{
		"app":      config.App(),
		"env":      config.GetEnv(),
		"instance": "test-stack",
	}

	tagMap := map[string]string{
		"app":          "happy/app",
		"env":          "happy/env",
		"instance":     "happy/instance",
		"owner":        "happy/meta/owner",
		"priority":     "happy/meta/priority",
		"slice":        "happy/meta/slice",
		"imagetag":     "happy/meta/imagetag",
		"imagetags":    "happy/meta/imagetags",
		"configsecret": "happy/meta/configsecret",
		"created":      "happy/meta/created-at",
		"updated":      "happy/meta/updated-at",
	}

	paramMap := map[string]string{
		"instance":     "stack_name",
		"slice":        "slice",
		"priority":     "priority",
		"imagetag":     "image_tag",
		"imagetags":    "image_tags",
		"configsecret": "happy_config_secret",
	}

	stackMeta := &StackMeta{
		stackName: "test-stack",
		DataMap:   dataMap,
		TagMap:    tagMap,
		paramMap:  paramMap,
	}

	// mock the backend
	ssm := mocks.NewMockSSMAPI(ctrl)
	retVal := "[\"stack_1\",\"stack_2\"]"
	ssm.EXPECT().GetParameter(gomock.Any()).Return(&retVal, nil)

	// mock the workspace GetTags method, used in setPriority()
	mockWorkspace1 := mocks.NewMockWorkspace(ctrl)
	mockWorkspace1.EXPECT().GetTags().Return(map[string]string{"tag-1": "testing-1"}, nil)
	mockWorkspace2 := mocks.NewMockWorkspace(ctrl)
	mockWorkspace2.EXPECT().GetTags().Return(map[string]string{"tag-2": "testing-2"}, nil)

	// mock the executor
	mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
	first := mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace1, nil)
	second := mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace2, nil)
	gomock.InOrder(first, second)

	backend, err := backend.NewAWSBackend(ctx, config, backend.WithSSMClient(ssm), backend.WithSecretsClient(secrets))
	r.NoError(err)

	stackMgr := NewStackService(backend, mockWorkspaceRepo)
	err = stackMeta.Update(ctx, "test-tag", make(map[string]string), "", stackMgr)
	r.NoError(err)
}
