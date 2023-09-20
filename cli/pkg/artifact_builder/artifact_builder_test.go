package artifact_builder

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "./testdata/test_config.yaml"
const testDockerComposePath = "./testdata/docker-compose.yml"

func TestCheckTagExists(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	workspace_repo.StartTFCWorkerPool(ctx)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	ecrApi := interfaces.NewMockECRAPI(ctrl)
	ecrApi.EXPECT().PutImage(gomock.Any(), gomock.Any()).Return(&ecr.PutImageOutput{}, nil).MaxTimes(3)
	ecrApi.EXPECT().GetAuthorizationToken(gomock.Any(), gomock.Any()).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []ecrtypes.AuthorizationData{
			{
				AuthorizationToken: aws.String("YTpiOmM6ZA=="),
				ProxyEndpoint:      aws.String("https://1234567.dkr.aws.czi.us-west-2.com"),
			},
		},
	}, nil).AnyTimes()
	ecrApi.EXPECT().BatchGetImage(gomock.Any(), gomock.Any()).Return(&ecr.BatchGetImageOutput{
		Images: []ecrtypes.Image{
			{
				ImageManifest: aws.String("manifest"),
			},
		},
	}, nil).MaxTimes(3)
	ecrApi.EXPECT().BatchGetRepositoryScanningConfiguration(gomock.Any(), gomock.Any()).Return(&ecr.BatchGetRepositoryScanningConfigurationOutput{
		ScanningConfigurations: []ecrtypes.RepositoryScanningConfiguration{
			{
				ScanOnPush: false,
			},
		},
	}, nil).AnyTimes()

	ecrApi.EXPECT().GetRegistryScanningConfiguration(gomock.Any(), gomock.Any()).Return(&ecr.GetRegistryScanningConfigurationOutput{
		RegistryId: aws.String("1234567"),
		ScanningConfiguration: &ecrtypes.RegistryScanningConfiguration{
			ScanType: ecrtypes.ScanTypeBasic,
		},
	}, nil).AnyTimes()

	buildConfig := NewBuilderConfig().WithBootstrap(bootstrapConfig).WithHappyConfig(happyConfig)
	buildConfig.Executor = util.NewDummyExecutor()
	backend, err := testbackend.NewBackend(ctx, ctrl, happyConfig.GetEnvironmentContext(), backend.WithECRClient(ecrApi), backend.WithExecutor(util.NewDummyExecutor()))
	r.NoError(err)

	configData, err := buildConfig.GetConfigData(ctx)
	r.NoError(err)
	configData.Services = make(map[string]ServiceConfig)
	configData.Services["frontend"] = ServiceConfig{
		Image:   "nginx",
		Build:   &ServiceBuild{},
		Network: map[string]interface{}{},
	}

	artifactBuilder := CreateArtifactBuilder(ctx).WithHappyConfig(happyConfig).WithConfig(buildConfig).WithBackend(backend)

	registryConfig := config.RegistryConfig{
		URL: "1234567.dkr.aws.czi.us-west-2.com/nginx",
	}
	serviceRegistries := backend.Conf().GetServiceRegistries()
	serviceRegistries["frontend"] = &registryConfig
	serviceRegistries["ecr_1"] = &registryConfig
	r.NotNil(serviceRegistries)
	r.Len(serviceRegistries, 2)

	_, err = artifactBuilder.CheckImageExists(ctx, "a")
	r.NoError(err)

	err = artifactBuilder.RetagImages(ctx, serviceRegistries, "latest", []string{"latest"}, []string{})
	r.NoError(err)

	err = artifactBuilder.RegistryLogin(context.Background())
	r.NoError(err)
	err = artifactBuilder.Build(ctx)
	r.NoError(err)
	err = artifactBuilder.Push(ctx, []string{"latest"})
	r.NoError(err)
}

func TestBuildAndPush(t *testing.T) {
	workspace_repo.StartTFCWorkerPool(context.Background())
	r := require.New(t)
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	// mock docker login
	dockerRegistry := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// TODO: assert the token is what we want
			fmt.Fprintln(w, "Hello, client")
		},
		),
	)
	defer dockerRegistry.Close()

	// mock ecr
	ecrApi := interfaces.NewMockECRAPI(ctrl)
	ecrApi.EXPECT().PutImage(gomock.Any(), gomock.Any()).Return(&ecr.PutImageOutput{}, nil).MaxTimes(3)
	ecrApi.EXPECT().GetAuthorizationToken(gomock.Any(), gomock.Any()).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []ecrtypes.AuthorizationData{
			{
				AuthorizationToken: aws.String("YTpiOmM6ZA=="),
				ProxyEndpoint:      aws.String(dockerRegistry.URL),
			},
		},
	}, nil).MaxTimes(2)
	ecrApi.EXPECT().BatchGetImage(gomock.Any(), gomock.Any()).Return(&ecr.BatchGetImageOutput{
		Images: []ecrtypes.Image{
			{
				ImageManifest: aws.String("manifest"),
			},
		},
	}, nil).MaxTimes(5)

	buildConfig := NewBuilderConfig().WithBootstrap(bootstrapConfig).WithHappyConfig(happyConfig)
	buildConfig.Executor = util.NewDummyExecutor()
	backend, err := testbackend.NewBackend(ctx, ctrl, happyConfig.GetEnvironmentContext(), backend.WithECRClient(ecrApi), backend.WithExecutor(util.NewDummyExecutor()))
	r.NoError(err)

	buildConfig.SetConfigData(&ConfigData{
		Services: map[string]ServiceConfig{"service1": {
			Image:   "nginx",
			Build:   &ServiceBuild{},
			Network: map[string]interface{}{},
		}},
	})
	artifactBuilder := CreateArtifactBuilder(ctx).WithHappyConfig(happyConfig).WithConfig(buildConfig)

	err = artifactBuilder.BuildAndPush(ctx)
	r.Error(err)

	artifactBuilder = artifactBuilder.WithBackend(backend)

	err = artifactBuilder.BuildAndPush(ctx)
	r.NoError(err)

	artifactBuilder = CreateArtifactBuilder(ctx).WithHappyConfig(happyConfig).WithConfig(buildConfig).WithBackend(backend).WithTags([]string{"test"})

	err = artifactBuilder.BuildAndPush(ctx)
	r.NoError(err)
}
