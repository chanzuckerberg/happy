package artifact_builder

import (
	"context"
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestNewBuilderConfig(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	bootstrap := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
	}
	happyConfig, err := config.NewHappyConfig(ctx, bootstrap)
	r.NoError(err)

	builderConfig := NewBuilderConfig(bootstrap, happyConfig)
	r.NotNil(builderConfig)

	_, err = builderConfig.GetContainers()
	r.NoError(err)
}

func TestNewBuilderConfigProfiles(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	bootstrap := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
	}

	happyConfig, err := config.NewHappyConfig(ctx, bootstrap)
	r.NoError(err)

	bc := NewBuilderConfig(bootstrap, happyConfig)
	r.NotNil(bc)

	// TODO: figure out why this is empty
	configData, err := bc.getConfigData()
	r.NoError(err)

	spew.Dump(configData)
	r.Nil(configData)

	_, err = bc.GetContainers()
	r.NoError(err)
}
