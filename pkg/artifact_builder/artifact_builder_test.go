package artifact_builder

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"
const testDockerComposePath = "../config/testdata/docker-compose.yml"

func TestCheckTagExists(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	secrets := mocks.NewMockSecretsManagerAPI(ctrl)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secrets.EXPECT().GetSecretValueWithContext(gomock.Any(), gomock.Any()).Return(&secretsmanager.GetSecretValueOutput{
		SecretBinary: []byte(testVal),
		SecretString: &testVal,
	},
		nil)

	stsApi := mocks.NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentityWithContext(gomock.Any(), gomock.Any()).Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil)

	happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
	r.NoError(err)

	buildConfig := NewBuilderConfig(bootstrapConfig, happyConfig)
	backend, err := backend.NewAWSBackend(ctx, happyConfig, backend.WithSecretsClient(secrets), backend.WithSTSClient(stsApi))
	r.NoError(err)

	artifactBuilder := NewArtifactBuilder(buildConfig, backend)

	serviceRegistries := backend.Conf().GetServiceRegistries()
	r.NotNil(serviceRegistries)
	r.True(len(serviceRegistries) > 0)

	_, err = artifactBuilder.CheckImageExists(serviceRegistries, "a")
	// Behind the scenes, an invocation of docker-compose is made, and it doesn't exist in github action image
	fmt.Printf("Error: %v\n", err)
	r.True(err == nil || strings.Contains(err.Error(), "executable file not found in $PATH") || strings.Contains(err.Error(), "process failure"))
}
