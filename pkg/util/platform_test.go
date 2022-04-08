package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainerPlatformParsing(t *testing.T) {
	r := require.New(t)

	r.NotEmpty(getUserContainerPlatform())
	r.NotEmpty(GetUserContainerPlatform())
	r.Equal("linux/amd64", GetSystemContainerPlatform("x86_64"))
	r.Equal("linux/amd64", GetSystemContainerPlatform("x86-64"))
	r.Equal("linux/arm64", GetSystemContainerPlatform("aarch64"))
	r.Equal("linux/amd64", GetSystemContainerPlatform("linux/amd64"))
	r.Equal("linux/arm64", GetSystemContainerPlatform("linux/arm64"))
	r.Equal("linux/arm64", GetSystemContainerPlatform("linux/arm64"))
}
