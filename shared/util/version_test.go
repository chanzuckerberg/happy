package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	r := require.New(t)

	release := GetVersion()
	r.NotNil(release)
	r.Equal("version: undefined\ngit_sha: undefined", release.String())
}
