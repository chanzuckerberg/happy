package stack_mgr

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	happyMocks "github.com/chanzuckerberg/happy/mocks"
	config "github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRemoveSucceed(t *testing.T) {
	r := require.New(t)
	mockCtrl := gomock.NewController(t)

	env := "rdev"
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

	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(mockCtrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	}, nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)
	r.NotNil(awsSecretMgr)

	for _, testCase := range testData {
		// TODO mock the config interfarce instead
		config, err := NewTestHappyConfig(t, testFilePath, env, awsSecretMgr)
		r.NoError(err)

		mockWorkspace := happyMocks.NewMockWorkspace(mockCtrl)
		mockWorkspace.EXPECT().Run(gomock.Any()).Return(nil)

		mockWorkspaceRepo := happyMocks.NewMockWorkspaceRepoIface(mockCtrl)
		mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace, nil)

		mockParamStore := happyMocks.NewMockParamStoreBackend(mockCtrl)
		testParamStoreData := testCase.input
		mockParamStore.EXPECT().GetParameter(gomock.Any()).Return(&testParamStoreData, nil)
		mockParamStore.EXPECT().AddParams("/happy/rdev/stacklist", testCase.expect).Return(nil)

		m := NewStackService(config, mockParamStore, mockWorkspaceRepo)

		err = m.Remove(testStackName)
		r.NoError(err)
	}
}

func TestAddSucceed(t *testing.T) {
	r := require.New(t)
	mockCtrl := gomock.NewController(t)

	env := "rdev"
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

	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(mockCtrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	}, nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)
	r.NotNil(awsSecretMgr)

	for _, testCase := range testData {
		config, err := NewTestHappyConfig(t, testFilePath, env, awsSecretMgr)
		r.NoError(err)

		mockWorkspace := happyMocks.NewMockWorkspace(mockCtrl)
		mockWorkspace.EXPECT().Run(gomock.Any()).Return(nil)
		mockWorkspace.EXPECT().Wait().Return(nil)

		mockWorkspaceRepo := happyMocks.NewMockWorkspaceRepoIface(mockCtrl)
		mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace, nil)
		// the second call of GetWorkspace occurs after the workspace creation,
		// for purpose of verifying that the workspace has indeed been created
		mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(mockWorkspace, nil)

		mockParamStore := happyMocks.NewMockParamStoreBackend(mockCtrl)
		testParamStoreData := testCase.input
		mockParamStore.EXPECT().GetParameter(gomock.Any()).Return(&testParamStoreData, nil)
		mockParamStore.EXPECT().AddParams("/happy/rdev/stacklist", testCase.expect).Return(nil)

		m := NewStackService(config, mockParamStore, mockWorkspaceRepo)

		_, err = m.Add(testStackName)
		r.NoError(err)
	}
}
