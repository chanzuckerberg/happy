package artifact_builder

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/testbackend"
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

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
	r.NoError(err)

	ecrApi := testbackend.NewMockECRAPI(ctrl)
	ecrApi.EXPECT().PutImage(gomock.Any()).Return(&ecr.PutImageOutput{}, nil).MaxTimes(3)
	ecrApi.EXPECT().GetAuthorizationTokenWithContext(gomock.Any(), gomock.Any()).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			{
				AuthorizationToken: aws.String("YTpiOmM6ZA=="),
				ProxyEndpoint:      aws.String("https://1234567.dkr.aws.czi.us-west-2.com"),
			},
		},
	}, nil)
	ecrApi.EXPECT().BatchGetImage(gomock.Any()).Return(&ecr.BatchGetImageOutput{
		Images: []*ecr.Image{
			{
				ImageManifest: aws.String("manifest"),
			},
		},
	}, nil).MaxTimes(3)

	buildConfig := NewBuilderConfig(bootstrapConfig, happyConfig).WithExecutor(NewDummyExecutor())
	backend, err := testbackend.NewBackend(ctx, ctrl, happyConfig, backend.WithECRClient(ecrApi))
	r.NoError(err)

	configData := buildConfig.GetConfigData()
	configData.Services = make(map[string]ServiceConfig)
	configData.Services["frontend"] = ServiceConfig{
		Image:   "nginx",
		Build:   &ServiceBuild{},
		Network: map[string]interface{}{},
	}

	artifactBuilder := NewArtifactBuilder(buildConfig, backend)

	registryConfig := config.RegistryConfig{
		Url: "1234567.dkr.aws.czi.us-west-2.com/nginx",
	}
	serviceRegistries := backend.Conf().GetServiceRegistries()
	serviceRegistries["frontend"] = &registryConfig
	serviceRegistries["ecr_1"] = &registryConfig
	r.NotNil(serviceRegistries)
	r.True(len(serviceRegistries) > 0)

	_, err = artifactBuilder.CheckImageExists(serviceRegistries, "a")
	// Behind the scenes, an invocation of docker-compose is made, and it doesn't exist in github action image
	fmt.Printf("Error: %v\n", err)
	r.True(err == nil || strings.Contains(err.Error(), "executable file not found in $PATH") || strings.Contains(err.Error(), "process failure"))

	err = artifactBuilder.RetagImages(serviceRegistries, map[string]string{"frontend": "foo"}, "latest", []string{"latest"}, []string{})
	r.NoError(err)

	err = artifactBuilder.RegistryLogin(context.Background())
	r.Error(err)
	err = artifactBuilder.Build()
	r.NoError(err)
	err = artifactBuilder.Push(serviceRegistries, map[string]string{"frontend": "foo"}, []string{"latest"})
	r.NoError(err)
}
