package stack

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/chanzuckerberg/happy/shared/workspace_repo"
)

func TestApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := require.New(t)
	ctx := context.Background()

	testStackMeta := &StackMeta{
		StackName: "test-stack",
		App:       "test-app",
		Env:       "rdev",
		Owner:     "test-owner",
	}
	// mock the workspace
	// NOTE SetVars is expected to be called 5 times
	// NOTE metaTags is generated from tagMap values mapped to dataMap values
	metaTags, err := json.Marshal(testStackMeta)
	r.NoError(err)

	testVersionId := "test_version_id"
	mockWorkspace1 := workspace_repo.NewMockWorkspace(ctrl)
	mockWorkspace1.EXPECT().SetVars(ctx, "happymeta_", string(metaTags), gomock.Any(), false).Return(nil)
	metaKeys := map[string]any{}
	err = json.Unmarshal(metaTags, &metaKeys)
	r.NoError(err)
	for k, v := range metaKeys {
		mockWorkspace1.EXPECT().SetVars(ctx, k, util.TagValueToString(v), gomock.Any(), false).Return(nil)
	}
	mockWorkspace1.EXPECT().GetTags(ctx).Return(map[string]string{}, nil).MaxTimes(2)
	mockWorkspace1.EXPECT().UploadVersion(ctx, gomock.Any()).Return(testVersionId, nil)
	mockWorkspace1.EXPECT().RunConfigVersion(ctx, testVersionId, gomock.Any()).Return(nil)
	mockWorkspace1.EXPECT().WaitWithOptions(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(2)

	stackService := NewMockStackServiceIface(ctrl)
	stackService.EXPECT().GetStackWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace1, nil).MaxTimes(2)
	stackService.EXPECT().NewStackMeta(gomock.Any()).Return(testStackMeta).MaxTimes(2)

	stack := NewStack(
		"test-stack",
		stackService,
	).WithMeta(testStackMeta)

	err = stack.Apply(ctx, "", options.WaitOptions{})
	r.NoError(err)

	err = stack.Wait(ctx, options.WaitOptions{})
	r.NoError(err)
}
