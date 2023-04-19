package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTFParser(t *testing.T) {
	services, err := ParseServices("./testdata/tf")
	r := require.New(t)
	r.NoError(err)
	r.True(len(services) > 0)
}
