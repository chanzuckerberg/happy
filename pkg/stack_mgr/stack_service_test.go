package stack_mgr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestWriteParam(t *testing.T) {
	testData := []struct {
		environment       string
		configPathValue   string
		expectedWritePath string
	}{
		{
			// uses the default when nothing is set
			environment:       "rdev",
			expectedWritePath: "/happy/rdev/stacklist",
		},
		{
			// uses the override when provided
			environment:       "rdev_with_stacklist_override",
			expectedWritePath: "/happy/api/rdev_with_stacklist_override/stacklist",
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
				Env:                     testCase.environment,
			}
			config, err := config.NewHappyConfig(bootstrapConfig)
			r.NoError(err)

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock))
			r.NoError(err)

			m := stack_mgr.NewStackService().WithBackend(backend)

			r.Equal(testCase.expectedWritePath, m.GetWritePath())
		})
	}
}

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
			config, err := config.NewHappyConfig(bootstrapConfig)
			r.NoError(err)

			mockWorkspace := mocks.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil)
			mockWorkspace.EXPECT().GetOutputs().Return(map[string]string{}, nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetLatestConfigVersionID().Return("123", nil).MaxTimes(100)
			mockWorkspace.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().Wait(gomock.Any(), gomock.Any()).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunStatus().Return("").MaxTimes(100)
			mockWorkspace.EXPECT().HasState(gomock.Any()).Return(true, nil).MaxTimes(100)
			mockWorkspace.EXPECT().RunConfigVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunID().Return("1234").MaxTimes(100)

			mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil).MaxTimes(100)

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock))
			r.NoError(err)

			m := stack_mgr.NewStackService().WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			err = m.Remove(ctx, testStackName, false)
			r.NoError(err)

			stacks, err := m.GetStacks(ctx)
			r.NoError(err)
			for _, stack := range stacks {
				_, err = stack.GetOutputs(ctx)
				r.NoError(err)
				stack.PrintOutputs(ctx)
				err = stack.PlanDestroy(ctx, false)
				r.NoError(err)
				r.Equal("", stack.GetStatus())
				hasState, err := m.HasState(ctx, stack.GetName())
				r.NoError(err)
				r.True(hasState)
			}
		})
	}
}

func TestRemoveWithLockSucceed(t *testing.T) {
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
			config, err := config.NewHappyConfig(bootstrapConfig)
			r.NoError(err)

			config.GetFeatures().EnableDynamoLocking = true

			mockWorkspace := mocks.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil)
			mockWorkspace.EXPECT().GetOutputs().Return(map[string]string{}, nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetLatestConfigVersionID().Return("123", nil).MaxTimes(100)
			mockWorkspace.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().Wait(gomock.Any(), gomock.Any()).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunStatus().Return("").MaxTimes(100)
			mockWorkspace.EXPECT().HasState(gomock.Any()).Return(true, nil).MaxTimes(100)
			mockWorkspace.EXPECT().RunConfigVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunID().Return("1234").MaxTimes(100)

			mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil).MaxTimes(100)

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			dynamoMock := interfaces.NewMockDynamoDB(ctrl)
			getItemRet := &dynamodb.GetItemOutput{}
			dynamoMock.EXPECT().GetItem(ctx, gomock.Any()).Return(getItemRet, nil)
			putItemRet := &dynamodb.PutItemOutput{}
			dynamoMock.EXPECT().PutItem(ctx, gomock.Any()).Return(putItemRet, nil)
			delItemRet := &dynamodb.DeleteItemOutput{}
			dynamoMock.EXPECT().DeleteItem(ctx, gomock.Any()).Return(delItemRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock), backend.WithDynamoDBClient(dynamoMock))
			r.NoError(err)

			m := stack_mgr.NewStackService().WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			err = m.Remove(ctx, testStackName, false)
			r.NoError(err)

			stacks, err := m.GetStacks(ctx)
			r.NoError(err)
			for _, stack := range stacks {
				_, err = stack.GetOutputs(ctx)
				r.NoError(err)
				stack.PrintOutputs(ctx)
				err = stack.PlanDestroy(ctx, false)
				r.NoError(err)
				r.Equal("", stack.GetStatus())
				hasState, err := m.HasState(ctx, stack.GetName())
				r.NoError(err)
				r.True(hasState)
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
			config, err := config.NewHappyConfig(bootstrapConfig)
			r.NoError(err)

			mockWorkspace := mocks.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil)
			mockWorkspace.EXPECT().Wait(gomock.Any(), gomock.Any()).Return(nil)

			mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil)
			// the second call of GetWorkspace occurs after the workspace creation,
			// for purpose of verifying that the workspace has indeed been created
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil)

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock))
			r.NoError(err)

			m := stack_mgr.NewStackService().WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			_, err = m.Add(ctx, testStackName, false)
			r.NoError(err)
		})
	}
}

func TestAddWithLockSucceed(t *testing.T) {
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
			config, err := config.NewHappyConfig(bootstrapConfig)
			r.NoError(err)

			config.GetFeatures().EnableDynamoLocking = true

			mockWorkspace := mocks.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil)
			mockWorkspace.EXPECT().Wait(gomock.Any(), gomock.Any()).Return(nil)

			mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil)
			// the second call of GetWorkspace occurs after the workspace creation,
			// for purpose of verifying that the workspace has indeed been created
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil)

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			dynamoMock := interfaces.NewMockDynamoDB(ctrl)
			getItemRet := &dynamodb.GetItemOutput{}
			dynamoMock.EXPECT().GetItem(ctx, gomock.Any()).Return(getItemRet, nil)
			putItemRet := &dynamodb.PutItemOutput{}
			dynamoMock.EXPECT().PutItem(ctx, gomock.Any()).Return(putItemRet, nil)
			delItemRet := &dynamodb.DeleteItemOutput{}
			dynamoMock.EXPECT().DeleteItem(ctx, gomock.Any()).Return(delItemRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config, backend.WithSSMClient(ssmMock), backend.WithDynamoDBClient(dynamoMock))
			r.NoError(err)

			m := stack_mgr.NewStackService().WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			_, err = m.Add(ctx, testStackName, false)
			r.NoError(err)
		})
	}
}
