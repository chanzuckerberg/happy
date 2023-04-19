package stack_mgr_test

import (
	"context"
	"testing"

	"github.com/chanzuckerberg/happy/cli/mocks"
	"github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := require.New(t)
	ctx := context.Background()

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}
	config, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	testStackMeta := &stack_mgr.StackMeta{
		StackName: "test-stack",
		App:       "test-app",
		Env:       "rdev",
	}
	// mock the workspace
	// NOTE SetVars is expected to be called 5 times
	// NOTE metaTags is generated from tagMap values mapped to dataMap values
	metaTags := "{\"happy/app\":\"test-app\",\"happy/env\":\"rdev\",\"happy/instance\":\"test-stack\",\"happy/meta/configsecret\":\"test-secret\"}"
	testVersionId := "test_version_id"
	mockWorkspace1 := mocks.NewMockWorkspace(ctrl)
	mockWorkspace1.EXPECT().SetVars(ctx, "happymeta_", metaTags, "Happy Path metadata", false).Return(nil)
	mockWorkspace1.EXPECT().GetTags(ctx).Return(map[string]string{}, nil).MaxTimes(2)
	mockWorkspace1.EXPECT().ResetCache().Return()
	mockWorkspace1.EXPECT().UploadVersion(ctx, gomock.Any(), gomock.Any()).Return(testVersionId, nil)
	mockWorkspace1.EXPECT().RunConfigVersion(ctx, testVersionId, gomock.Any()).Return(nil)
	mockWorkspace1.EXPECT().WaitWithOptions(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(2)

	stackService := mocks.NewMockStackServiceIface(ctrl)
	stackService.EXPECT().GetStackWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace1, nil)
	stackService.EXPECT().NewStackMeta(gomock.Any()).Return(testStackMeta).MaxTimes(2)
	stackService.EXPECT().GetConfig().Return(config).MaxTimes(2)

	mockDirProcessor := mocks.NewMockDirProcessor(ctrl)
	mockDirProcessor.EXPECT().Tarzip(gomock.Any(), gomock.Any()).Return(nil)

	stack := stack_mgr.NewStack(
		"test-stack",
		stackService,
		mockDirProcessor,
	)

	err = stack.Apply(ctx, options.WaitOptions{}, false)
	r.NoError(err)

	err = stack.Wait(ctx, options.WaitOptions{}, false)
	r.NoError(err)

	stack = stack.WithMeta(nil)
	r.NoError(err)
}
