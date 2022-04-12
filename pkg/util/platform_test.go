package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainerPlatformParsing(t *testing.T) {
	r := require.New(t)

	r.NotEmpty(getUserContainerPlatform())
	r.NotEmpty(GetUserContainerPlatform())

	sourceArch := []string{"x86_64", "x86-64", "aarch64", "linux/amd64", "linux/arm64"}
	targetArch := []string{"linux/amd64", "linux/amd64", "linux/arm64", "linux/amd64", "linux/arm64"}

	for index, arch := range sourceArch {
		plat, err := GetSystemContainerPlatform(arch)
		r.NoError(err)
		r.Equal(targetArch[index], plat)
	}
}
