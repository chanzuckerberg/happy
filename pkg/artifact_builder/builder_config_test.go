package artifact_builder

import (
	"context"
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNewBuilderConfig(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	bootstrap := &config.Bootstrap{HappyConfigPath: testFilePath}
	happyConfig, err := config.NewHappyConfig(ctx, bootstrap)
	r.NoError(err)

	builderConfig := NewBuilderConfig(bootstrap, happyConfig)
	r.NotNil(builderConfig)

	_, err = builderConfig.GetContainers()
	r.Error(err)
}
