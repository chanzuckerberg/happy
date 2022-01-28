package stack_mgr

import (
	"testing"

	happyMocks "github.com/chanzuckerberg/happy/mocks"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	env := "rdev"

	testFilePath := "../config/testdata/test_config.yaml"

	r := require.New(t)
	config, err := config.NewHappyConfig(testFilePath, env)
	r.NoError(err)

	dataMap := map[string]string{
		"app":      config.App(),
		"env":      config.DefaultEnv(),
		"instance": "test-stack",
	}

	tagMap := map[string]string{
		"app":          "happy/app",
		"env":          "happy/env",
		"instance":     "happy/instance",
		"owner":        "happy/meta/owner",
		"priority":     "happy/meta/priority",
		"imagetag":     "happy/meta/imagetag",
		"configsecret": "happy/meta/configsecret",
		"created":      "happy/meta/created-at",
		"updated":      "happy/meta/updated-at",
	}

	paramMap := map[string]string{
		"instance":     "stack_name",
		"priority":     "priority",
		"imagetag":     "image_tag",
		"configsecret": "happy_config_secret",
	}

	stackMeta := &StackMeta{
		stackName: "test-stack",
		DataMap:   dataMap,
		TagMap:    tagMap,
		paramMap:  paramMap,
	}

	// mock the backend
	mockCtrl := gomock.NewController(t)
	mockBackend := happyMocks.NewMockParamStoreBackend(mockCtrl)
	retVal := "[\"stack_1\",\"stack_2\"]"
	mockBackend.EXPECT().GetParameter(gomock.Any()).Return(&retVal, nil)

	// mock the workspace GetTags method, used in setPriority()
	mockWorkspace1 := happyMocks.NewMockWorkspace(mockCtrl)
	mockWorkspace1.EXPECT().GetTags().Return(map[string]string{"tag-1": "testing-1"}, nil)
	mockWorkspace2 := happyMocks.NewMockWorkspace(mockCtrl)
	mockWorkspace2.EXPECT().GetTags().Return(map[string]string{"tag-2": "testing-2"}, nil)

	// mock the executor
	mockWorkspaceRepo := happyMocks.NewMockWorkspaceRepoIface(mockCtrl)
	first := mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace1, nil)
	second := mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace2, nil)
	gomock.InOrder(first, second)

	stackMgr := NewStackService(config, mockBackend, mockWorkspaceRepo)
	err = stackMeta.Update("test-tag", stackMgr)
	r.NoError(err)
}

func TestGetTags(t *testing.T) {

}

func TestGetParameters(t *testing.T) {

}

func TestLoad(t *testing.T) {

}
