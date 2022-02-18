package stack_mgr_test

import (
	"context"
	"testing"

	"github.com/chanzuckerberg/happy/mocks"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/options"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	r := require.New(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}
	config, err := config.NewHappyConfig(ctx, bootstrapConfig)
	r.NoError(err)

	// StackMeta
	dataMap := map[string]string{
		"app":      "test-app",
		"env":      "rdev",
		"instance": "test-stack",
	}

	tagMap := map[string]string{
		"app":          "happy/app",
		"env":          "happy/env",
		"instance":     "happy/instance",
		"configsecret": "happy/meta/configsecret",
	}

	paramMap := map[string]string{
		"instance":     "stack_name",
		"priority":     "priority",
		"imagetag":     "image_tag",
		"configsecret": "happy_config_secret",
	}

	testStackMeta := &stack_mgr.StackMeta{
		StackName: "test-stack",
		DataMap:   dataMap,
		TagMap:    tagMap,
		ParamMap:  paramMap,
	}
	err = testStackMeta.Load(map[string]string{"happy/meta/configsecret": "test-secret"})
	r.NoError(err)

	// mock the workspace
	// NOTE SetVars is expected to be called 5 times
	// NOTE metaTags is generated from tagMap values mapped to dataMap values
	metaTags := "{\"happy/app\":\"test-app\",\"happy/env\":\"rdev\",\"happy/instance\":\"test-stack\",\"happy/meta/configsecret\":\"test-secret\"}"
	testVersionId := "test_version_id"
	mockWorkspace1 := mocks.NewMockWorkspace(ctrl)
	mockWorkspace1.EXPECT().SetVars("happymeta_", metaTags, "Happy Path metadata", false).Return(nil)
	for i := 0; i < len(paramMap); i++ {
		mockWorkspace1.EXPECT().SetVars(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	}
	mockWorkspace1.EXPECT().GetTags().Return(map[string]string{}, nil)
	mockWorkspace1.EXPECT().ResetCache().Return()
	mockWorkspace1.EXPECT().UploadVersion(gomock.Any()).Return(testVersionId, nil)
	mockWorkspace1.EXPECT().RunConfigVersion(testVersionId, gomock.Any()).Return(nil)
	mockWorkspace1.EXPECT().WaitWithOptions(gomock.Any()).Return(nil)

	stackService := mocks.NewMockStackServiceIface(ctrl)
	stackService.EXPECT().GetStackWorkspace(gomock.Any()).Return(mockWorkspace1, nil)
	stackService.EXPECT().NewStackMeta(gomock.Any()).Return(testStackMeta)
	stackService.EXPECT().GetConfig().Return(config).MaxTimes(2)

	mockDirProcessor := mocks.NewMockDirProcessor(ctrl)
	mockDirProcessor.EXPECT().Tarzip(gomock.Any(), gomock.Any()).Return(nil)

	stack := stack_mgr.NewStack(
		"test-stack",
		stackService,
		mockDirProcessor,
	)

	err = stack.Apply(options.WaitOptions{})
	r.NoError(err)
}
