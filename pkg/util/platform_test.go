package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlatformParsing(t *testing.T) {
	r := require.New(t)

	r.NotEmpty(getUserPlatform())
	r.NotEmpty(GetUserPlatform())
	r.Equal("linux/amd64", GetSystemPlatform("x86_64"))
	r.Equal("linux/arm64", GetSystemPlatform("aarch64"))
}
