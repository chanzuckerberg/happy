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

func TestTFParserOutputs(t *testing.T) {
	outputs, err := NewTfParser().ParseOutputs("./testdata")
	r := require.New(t)
	r.NoError(err)
	r.True(len(outputs) > 0)
}

func TestTFParserVariables(t *testing.T) {
	variables, err := NewTfParser().ParseVariables("./testdata")
	r := require.New(t)
	r.NoError(err)
	r.True(len(variables) > 0)
}
