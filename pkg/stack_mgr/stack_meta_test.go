package stack_mgr_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"
const testDockerComposePath = "../config/testdata/docker-compose.yml"

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	r := require.New(t)
	ctrl := gomock.NewController(t)

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

	stackMeta := &stack_mgr.StackMeta{
		StackName: "test-stack",
		DataMap:   dataMap,
		TagMap:    tagMap,
		ParamMap:  paramMap,
	}

	// mock the backend
	ssmMock := testbackend.NewMockSSMAPI(ctrl)
	retVal := "[\"stack_1\",\"stack_2\"]"
	ret := &ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{Value: &retVal},
	}
	ssmMock.EXPECT().GetParameterWithContext(gomock.Any(), gomock.Any()).Return(ret, nil)

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

	backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock))
	r.NoError(err)

	stackMgr := stack_mgr.NewStackService(backend, mockWorkspaceRepo)
	err = stackMeta.Update(ctx, "test-tag", make(map[string]string), "", stackMgr)
	r.NoError(err)
	r.Equal("{}", stackMeta.GetTags()["happy/meta/imagetags"])
	err = stackMeta.Update(ctx, "test-tag", map[string]string{"foo": "bar"}, "", stackMgr)
	r.NoError(err)
	r.Equal("{\"foo\":\"bar\"}", stackMeta.GetTags()["happy/meta/imagetags"])

}
