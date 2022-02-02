package artifact_builder

import (
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNewBuilderConfig(t *testing.T) {
	r := require.New(t)

	r.NotNil(NewBuilderConfig(&config.Bootstrap{DockerComposeConfigPath: "foobar"}, "bar"))
}
