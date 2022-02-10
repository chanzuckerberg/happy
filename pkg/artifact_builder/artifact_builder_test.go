package artifact_builder

import (
	"context"
	"fmt"
	"strings"
	"testing"

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

	buildConfig := NewBuilderConfig(bootstrapConfig, happyConfig)
	backend, err := testbackend.NewBackend(ctx, ctrl, happyConfig)
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
