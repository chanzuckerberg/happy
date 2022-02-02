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

	testVal := "{\"cluster_arn\":\"test_arn\",\"ecrs\":{\"ecr_1\":{\"url\":\"test_url_1\"}}}"
	mock.EXPECT().GetSecretValue(gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: &testVal,
	},
		nil)

	awsSecretMgr := config.GetAwsSecretMgrWithClient(mock)
	r.NotNil(awsSecretMgr)

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	happyConfig.SetSecretsBackend(awsSecretMgr)

	buildConfig := NewBuilderConfig(bootstrapConfig, "")
	artifactBuilder := NewArtifactBuilder(buildConfig, happyConfig)

	serviceRegistries, err := happyConfig.GetRdevServiceRegistries()
	r.NoError(err)

	imageExists, err := artifactBuilder.CheckImageExists(serviceRegistries, "a")
	// TODO(el): assert error is what we expect it to be
	r.Error(err)
	r.False(imageExists)
}
