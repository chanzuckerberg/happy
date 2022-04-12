package util

import (
	"github.com/containerd/containerd/platforms"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// All docker runtimes run images of linux/amd64 or linux/arm64 architecture -- even on darwin or windows.
// To see all supported runtimes on your machine, run: docker buildx ls
func normalizeContainerPlatfrom(platform v1.Platform) v1.Platform {
	platform.OS = "linux"
	return platforms.Normalize(platform)
}

func getUserContainerPlatform() v1.Platform {
	platform := platforms.DefaultSpec()
	return normalizeContainerPlatfrom(platform)
}

func GetSystemContainerPlatform(architecture string) (string, error) {
	platform, err := platforms.Parse(architecture)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse architecture")
	}
	return platforms.Format(normalizeContainerPlatfrom(platform)), nil
}

func GetUserContainerPlatform() string {
	return platforms.Format(getUserContainerPlatform())
}
