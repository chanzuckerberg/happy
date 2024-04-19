package stack

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type MockOrchestrator struct {
}

func (m *MockOrchestrator) GetEvents(ctx context.Context, stack string, services []string) error {
	return nil
}

func (m *MockOrchestrator) PrintLogs(ctx context.Context, stack string, services []string) error {
	return nil
}

func TestRemoveSucceed(t *testing.T) {
	testStackName := "test_stack"

	testData := []struct {
		input  string
		expect string
	}{
		{
			input:  fmt.Sprintf("[\"stack_1\",\"stack_2\",\"%s\"]", testStackName),
			expect: "[\"stack_1\",\"stack_2\"]",
		},
		{
			input:  fmt.Sprintf("[\"%s\"]", testStackName),
			expect: "[]",
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

			mockWorkspace := workspace_repo.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(ctx).Return(nil)
			mockWorkspace.EXPECT().GetOutputs(ctx).Return(map[string]string{}, nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetLatestConfigVersionID(ctx).Return("123", nil).MaxTimes(100)
			mockWorkspace.EXPECT().Run(ctx).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().Wait(gomock.Any()).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunStatus(ctx).Return("").MaxTimes(100)
			mockWorkspace.EXPECT().HasState(gomock.Any()).Return(true, nil).MaxTimes(100)
			mockWorkspace.EXPECT().RunConfigVersion(ctx, gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunID().Return("1234").MaxTimes(100)
			mockWorkspace.EXPECT().SetVars(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().UploadVersion(ctx, gomock.Any()).Return("123", nil).MaxTimes(100)
			mockWorkspace.EXPECT().WaitWithOptions(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspaceRepo := workspace_repo.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil).MaxTimes(100)

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(ssmRet, nil).AnyTimes()
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil).Times(2)

			dynamoMock := interfaces.NewMockDynamoDB(ctrl)
			getItemRet := &dynamodb.GetItemOutput{}
			dynamoMock.EXPECT().GetItem(ctx, gomock.Any()).Return(getItemRet, nil)
			putItemRet := &dynamodb.PutItemOutput{}
			dynamoMock.EXPECT().PutItem(ctx, gomock.Any()).Return(putItemRet, nil)
			delItemRet := &dynamodb.DeleteItemOutput{}
			dynamoMock.EXPECT().DeleteItem(ctx, gomock.Any()).Return(delItemRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config.GetEnvironmentContext(), backend.WithSSMClient(ssmMock), backend.WithDynamoDBClient(dynamoMock))
			r.NoError(err)

			m := NewStackService(config.GetEnv(), config.App()).WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			err = m.Remove(ctx, testStackName)
			r.NoError(err)

			stacks, err := m.GetStacks(ctx)
			r.NoError(err)
			for _, stack := range stacks {
				_, err = stack.GetOutputs(ctx)
				r.NoError(err)
				stack.PrintOutputs(ctx)
				waitoptions := options.WaitOptions{
					StackName:    testStackName,
					Orchestrator: &MockOrchestrator{},
					Services:     config.GetServices(),
				}
				err = stack.Destroy(ctx, "", waitoptions)
				r.NoError(err)
				r.Equal("", stack.GetStatus(ctx))
				hasState, err := m.HasState(ctx, stack.Name)
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

			mockWorkspace := workspace_repo.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(ctx).Return(nil)
			mockWorkspace.EXPECT().GetOutputs(ctx).Return(map[string]string{}, nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetLatestConfigVersionID(ctx).Return("123", nil).MaxTimes(100)
			mockWorkspace.EXPECT().Run(ctx, gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().Wait(gomock.Any()).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunStatus(ctx).Return("").MaxTimes(100)
			mockWorkspace.EXPECT().HasState(gomock.Any()).Return(true, nil).MaxTimes(100)
			mockWorkspace.EXPECT().RunConfigVersion(ctx, gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().GetCurrentRunID().Return("1234").MaxTimes(100)
			mockWorkspace.EXPECT().SetVars(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspace.EXPECT().UploadVersion(ctx, gomock.Any()).Return("123", nil).MaxTimes(100)
			mockWorkspace.EXPECT().WaitWithOptions(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(100)
			mockWorkspaceRepo := workspace_repo.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(mockWorkspace, nil).MaxTimes(100)

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(ssmRet, nil).AnyTimes()
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil).Times(2)

			dynamoMock := interfaces.NewMockDynamoDB(ctrl)
			getItemRet := &dynamodb.GetItemOutput{}
			dynamoMock.EXPECT().GetItem(ctx, gomock.Any()).Return(getItemRet, nil)
			putItemRet := &dynamodb.PutItemOutput{}
			dynamoMock.EXPECT().PutItem(ctx, gomock.Any()).Return(putItemRet, nil)
			delItemRet := &dynamodb.DeleteItemOutput{}
			dynamoMock.EXPECT().DeleteItem(ctx, gomock.Any()).Return(delItemRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config.GetEnvironmentContext(), backend.WithSSMClient(ssmMock), backend.WithDynamoDBClient(dynamoMock))
			r.NoError(err)

			m := NewStackService(config.GetEnv(), config.App()).WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			err = m.Remove(ctx, testStackName)
			r.NoError(err)

			stacks, err := m.GetStacks(ctx)
			r.NoError(err)
			for _, stack := range stacks {
				_, err = stack.GetOutputs(ctx)
				r.NoError(err)
				stack.PrintOutputs(ctx)
				waitoptions := options.WaitOptions{
					StackName:    testStackName,
					Orchestrator: &MockOrchestrator{},
					Services:     config.GetServices(),
				}
				err = stack.Destroy(ctx, "", waitoptions)
				r.NoError(err)
				r.Equal("", stack.GetStatus(ctx))
				hasState, err := m.HasState(ctx, stack.Name)
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

			mockWorkspace := workspace_repo.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(ctx).Return(nil)
			mockWorkspace.EXPECT().Wait(gomock.Any()).Return(nil)

			mockWorkspaceRepo := workspace_repo.NewMockWorkspaceRepoIface(ctrl)
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
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil).Times(2)

			dynamoMock := interfaces.NewMockDynamoDB(ctrl)
			getItemRet := &dynamodb.GetItemOutput{}
			dynamoMock.EXPECT().GetItem(ctx, gomock.Any()).Return(getItemRet, nil)
			putItemRet := &dynamodb.PutItemOutput{}
			dynamoMock.EXPECT().PutItem(ctx, gomock.Any()).Return(putItemRet, nil)
			delItemRet := &dynamodb.DeleteItemOutput{}
			dynamoMock.EXPECT().DeleteItem(ctx, gomock.Any()).Return(delItemRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config.GetEnvironmentContext(), backend.WithSSMClient(ssmMock), backend.WithDynamoDBClient(dynamoMock))
			r.NoError(err)

			m := NewStackService(config.GetEnv(), config.App()).WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			_, err = m.Add(ctx, testStackName)
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

			mockWorkspace := workspace_repo.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(ctx, gomock.Any()).Return(nil)
			mockWorkspace.EXPECT().Wait(gomock.Any()).Return(nil)

			mockWorkspaceRepo := workspace_repo.NewMockWorkspaceRepoIface(ctrl)
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
			ssmMock.EXPECT().PutParameter(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil).Times(2)

			dynamoMock := interfaces.NewMockDynamoDB(ctrl)
			getItemRet := &dynamodb.GetItemOutput{}
			dynamoMock.EXPECT().GetItem(ctx, gomock.Any()).Return(getItemRet, nil)
			putItemRet := &dynamodb.PutItemOutput{}
			dynamoMock.EXPECT().PutItem(ctx, gomock.Any()).Return(putItemRet, nil)
			delItemRet := &dynamodb.DeleteItemOutput{}
			dynamoMock.EXPECT().DeleteItem(ctx, gomock.Any()).Return(delItemRet, nil)

			backend, err := testbackend.NewBackend(ctx, ctrl, config.GetEnvironmentContext(), backend.WithSSMClient(ssmMock), backend.WithDynamoDBClient(dynamoMock))
			r.NoError(err)

			m := NewStackService(config.GetEnv(), config.App()).WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			_, err = m.Add(ctx, testStackName)
			r.NoError(err)
		})
	}
}

func TestGetStacksSucceed(t *testing.T) {
	testData := []struct {
		input                 string
		expect                []string
		namespacedParamExists bool
	}{
		{
			input:                 "[\"stack_1\",\"stack_2\"]",
			expect:                []string{"stack_1", "stack_2"},
			namespacedParamExists: false,
		},
		{
			input:                 "[\"stack_a\",\"stack_b\",\"stack_c\"]",
			expect:                []string{"stack_a", "stack_b", "stack_c"},
			namespacedParamExists: true,
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

			ssmMock := interfaces.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{Value: &testParamStoreData},
			}

			if testCase.namespacedParamExists {
				ssmMock.EXPECT().GetParameter(gomock.Any(), &ssm.GetParameterInput{Name: aws.String("/happy/test-app/rdev/stacklist")}).Return(ssmRet, nil)
			} else {
				ssmMock.EXPECT().GetParameter(gomock.Any(), &ssm.GetParameterInput{Name: aws.String("/happy/test-app/rdev/stacklist")}).Return(nil, errors.New("ParameterNotFound"))
				ssmMock.EXPECT().GetParameter(gomock.Any(), &ssm.GetParameterInput{Name: aws.String("/happy/rdev/stacklist")}).Return(ssmRet, nil)
			}

			backend, err := testbackend.NewBackend(ctx, ctrl, config.GetEnvironmentContext(), backend.WithSSMClient(ssmMock))
			r.NoError(err)

			mockWorkspaceRepo := workspace_repo.NewMockWorkspaceRepoIface(ctrl)
			m := NewStackService(config.GetEnv(), config.App()).WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)

			stacks, err := m.GetStacks(ctx)
			r.NoError(err)
			stackNames := []string{}
			for _, stack := range stacks {
				stackNames = append(stackNames, stack.Name)
			}

			r.ElementsMatch(testCase.expect, stackNames)
		})
	}
}
