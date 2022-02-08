package stack_mgr

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	happyMocks "github.com/chanzuckerberg/happy/mocks"
	config "github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/options"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	env := "rdev"

	r := require.New(t)
	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(ctrl)

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	}, nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)
	r.NotNil(awsSecretMgr)

	config, err := NewTestHappyConfig(t, testFilePath, env, awsSecretMgr)
	r.NoError(err)

	mockCtrl := gomock.NewController(t)

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

	testStackMeta := &StackMeta{
		stackName: "test-stack",
		DataMap:   dataMap,
		TagMap:    tagMap,
		paramMap:  paramMap,
	}
	err = testStackMeta.Load(map[string]string{"happy/meta/configsecret": "test-secret"})
	r.NoError(err)

	// mock the workspace
	// NOTE SetVars is expected to be called 5 times
	// NOTE metaTags is generated from tagMap values mapped to dataMap values
	metaTags := "{\"happy/app\":\"test-app\",\"happy/env\":\"rdev\",\"happy/instance\":\"test-stack\",\"happy/meta/configsecret\":\"test-secret\"}"
	testVersionId := "test_version_id"
	mockWorkspace1 := happyMocks.NewMockWorkspace(mockCtrl)
	mockWorkspace1.EXPECT().SetVars("happymeta_", metaTags, "Happy Path metadata", false).Return(nil)
	for i := 0; i < len(paramMap); i++ {
		mockWorkspace1.EXPECT().SetVars(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	}
	mockWorkspace1.EXPECT().GetTags().Return(map[string]string{}, nil)
	mockWorkspace1.EXPECT().ResetCache().Return()
	mockWorkspace1.EXPECT().UploadVersion(gomock.Any()).Return(testVersionId, nil)
	mockWorkspace1.EXPECT().RunConfigVersion(testVersionId, gomock.Any()).Return(nil)
	mockWorkspace1.EXPECT().WaitWithOptions(gomock.Any()).Return(nil)

	stackService := NewMockStackServiceIface(mockCtrl)
	stackService.EXPECT().GetStackWorkspace(gomock.Any()).Return(mockWorkspace1, nil)
	stackService.EXPECT().NewStackMeta(gomock.Any()).Return(testStackMeta)
	stackService.EXPECT().GetConfig().Return(config)

	mockDirProcessor := happyMocks.NewMockDirProcessor(mockCtrl)
	mockDirProcessor.EXPECT().Tarzip(gomock.Any(), gomock.Any()).Return(nil)

	stack := &Stack{
		stackService: stackService,
		stackName:    "test-stack",
		dirProcessor: mockDirProcessor,
	}

	err = stack.Apply(options.WaitOptions{})
	r.NoError(err)
}
