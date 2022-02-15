package testbackend

import (
	"context"
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../../../config/testdata/test_config.yaml"
const testDockerComposePath = "../../../config/testdata/docker-compose.yml"

func TestAWSBackend(t *testing.T) {
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

	_, err = NewBackend(ctx, ctrl, happyConfig)
	r.NoError(err)
}
