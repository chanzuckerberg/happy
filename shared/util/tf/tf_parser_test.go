package tf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTFParserServices(t *testing.T) {
	services, err := NewTfParser().ParseServices("./testdata")
	r := require.New(t)
	r.NoError(err)
	r.True(len(services) > 0)
}
