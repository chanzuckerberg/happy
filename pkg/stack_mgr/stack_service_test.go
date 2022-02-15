package stack_mgr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/testbackend"
	config "github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRemoveSucceed(t *testing.T) {
	testStackName := "test_stack"

	testData := []struct {
		input  string
		expect string
	}{
		{
			fmt.Sprintf("[\"stack_1\",\"stack_2\",\"%s\"]", testStackName),
			"[\"stack_1\",\"stack_2\"]",
		},
		{
			fmt.Sprintf("[\"%s\"]", testStackName),
			"[]",
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			bootstrapConfig := &config.Bootstrap{
				HappyConfigPath:         testFilePath,
				DockerComposeConfigPath: testDockerComposePath,
				Env:                     "rdev",
			}
			config, err := config.NewHappyConfig(ctx, bootstrapConfig)
			r.NoError(err)

			mockWorkspace := mocks.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(gomock.Any()).Return(nil)
			mockWorkspace.EXPECT().GetOutputs().Return(map[string]string{}, nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetLatestConfigVersionID().Return("", nil).MaxTimes(100)

			mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace, nil).MaxTimes(100)

			ssmMock := testbackend.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock))
			r.NoError(err)

			m := stack_mgr.NewStackService(backend, mockWorkspaceRepo)

			err = m.Remove(ctx, testStackName)
			r.NoError(err)

			stacks, err := m.GetStacks(ctx)
			r.NoError(err)
			for _, stack := range stacks {
				_, err = stack.GetOutputs()
				r.NoError(err)
				stack.PrintOutputs()
				stack.Destroy()
			}
		})
	}
}

func TestAddSucceed(t *testing.T) {
	testStackName := "test_stack"

	testData := []struct {
		input  string
		expect string
	}{
		{
			"[\"stack_1\",\"stack_2\"]",
			fmt.Sprintf("[\"stack_1\",\"stack_2\",\"%s\"]", testStackName),
		},
		{
			"[]",
			fmt.Sprintf("[\"%s\"]", testStackName),
		},
		{
			fmt.Sprintf("[\"%s\"]", testStackName),
			fmt.Sprintf("[\"%s\"]", testStackName),
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			bootstrapConfig := &config.Bootstrap{
				HappyConfigPath:         testFilePath,
				DockerComposeConfigPath: testDockerComposePath,
				Env:                     "rdev",
			}
			config, err := config.NewHappyConfig(ctx, bootstrapConfig)
			r.NoError(err)

			mockWorkspace := mocks.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(gomock.Any()).Return(nil)
			mockWorkspace.EXPECT().Wait().Return(nil)

			mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace, nil)
			// the second call of GetWorkspace occurs after the workspace creation,
			// for purpose of verifying that the workspace has indeed been created
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace, nil)

			ssmMock := testbackend.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock))
			r.NoError(err)

			m := stack_mgr.NewStackService(backend, mockWorkspaceRepo)

			_, err = m.Add(ctx, testStackName)
			r.NoError(err)
		})
	}
}
