package artifact_builder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBuilderConfig(t *testing.T) {
	r := require.New(t)

	r.NotNil(NewBuilderConfig("foo", "bar"))
}
