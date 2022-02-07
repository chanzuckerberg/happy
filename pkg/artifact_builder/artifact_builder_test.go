package artifact_builder

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	cziAWS "github.com/chanzuckerberg/go-misc/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"
const testDockerComposePath = "../config/testdata/docker-compose.yml"

func TestCheckTagExists(t *testing.T) {
	r := require.New(t)
	ctrl := gomock.NewController(t)
	client := cziAWS.Client{}
	_, mock := client.WithMockSecretsManager(ctrl)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	},
		nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)
	r.NotNil(awsSecretMgr)

	happyConfig, err := config.NewHappyConfigWithSecretsBackend(bootstrapConfig, awsSecretMgr)
	r.NoError(err)

	buildConfig := NewBuilderConfig(bootstrapConfig, "", happyConfig.GetDockerRepo())
	artifactBuilder := NewArtifactBuilder(buildConfig, happyConfig)

	serviceRegistries := happyConfig.GetRdevServiceRegistries()
	r.NotNil(serviceRegistries)
	r.True(len(serviceRegistries) > 0)

	imageExists, err := artifactBuilder.CheckImageExists(serviceRegistries, "a")
	// TODO(el): assert error is what we expect it to be
	r.NoError(err)
	r.True(imageExists)
}
