package stack_mgr

import (
	"testing"

	config "github.com/chanzuckerberg/happy/shared/config"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"
const testDockerComposePath = "../config/testdata/docker-compose.yml"

func TestUpdate(t *testing.T) {
	r := require.New(t)
	options := NewStackManagementOptions("stack1")
	r.Equal(options.StackName, "stack1")

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	config, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	options = options.WithHappyConfig(config)
	r.NotNil(options.HappyConfig)
}
