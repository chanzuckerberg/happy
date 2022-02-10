package stack_mgr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
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

			secrets := mocks.NewMockSecretsManagerAPI(ctrl)
			testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
			secrets.EXPECT().GetSecretValueWithContext(gomock.Any(), gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
				SecretBinary: []byte(testVal),
				SecretString: &testVal,
			}, nil)

			stsApi := mocks.NewMockSTSAPI(ctrl)
			stsApi.EXPECT().GetCallerIdentityWithContext(gomock.Any(), gomock.Any()).Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil)

			bootstrapConfig := &config.Bootstrap{
				HappyConfigPath:         testFilePath,
				DockerComposeConfigPath: testDockerComposePath,
				Env:                     "rdev",
			}
			config, err := config.NewHappyConfig(ctx, bootstrapConfig)
			r.NoError(err)

			mockWorkspace := mocks.NewMockWorkspace(ctrl)
			mockWorkspace.EXPECT().Run(gomock.Any()).Return(nil)

			mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
			mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace, nil)

			ssmMock := mocks.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			backend, err := backend.NewAWSBackend(ctx, config, backend.WithSSMClient(ssmMock), backend.WithSecretsClient(secrets), backend.WithSTSClient(stsApi))
			r.NoError(err)

			m := stack_mgr.NewStackService(backend, mockWorkspaceRepo)

			err = m.Remove(ctx, testStackName)
			r.NoError(err)
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

			secrets := mocks.NewMockSecretsManagerAPI(ctrl)
			testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
			secrets.EXPECT().GetSecretValueWithContext(gomock.Any(), gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
				SecretBinary: []byte(testVal),
				SecretString: &testVal,
			}, nil)

			stsApi := mocks.NewMockSTSAPI(ctrl)
			stsApi.EXPECT().GetCallerIdentityWithContext(gomock.Any(), gomock.Any()).Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil)

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

			ssmMock := mocks.NewMockSSMAPI(ctrl)
			testParamStoreData := testCase.input
			ssmRet := &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{Value: &testParamStoreData},
			}

			ssmPutRet := &ssm.PutParameterOutput{}
			ssmMock.EXPECT().GetParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmRet, nil)
			ssmMock.EXPECT().PutParameterWithContext(gomock.Any(), gomock.Any()).Return(ssmPutRet, nil)

			backend, err := backend.NewAWSBackend(ctx, config, backend.WithSSMClient(ssmMock), backend.WithSecretsClient(secrets), backend.WithSTSClient(stsApi))
			r.NoError(err)

			m := stack_mgr.NewStackService(backend, mockWorkspaceRepo)

			_, err = m.Add(ctx, testStackName)
			r.NoError(err)
		})
	}
}
