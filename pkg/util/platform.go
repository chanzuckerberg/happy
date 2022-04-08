package util

import (
	"fmt"

	"github.com/containerd/containerd/platforms"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func normalizeContainerPlatfrom(platform v1.Platform) v1.Platform {
	platform.OS = "linux"
	return platforms.Normalize(platform)
}

func getUserContainerPlatform() v1.Platform {
	platform := platforms.DefaultSpec()
	return normalizeContainerPlatfrom(platform)
}

func GetSystemContainerPlatform(architecture string) string {
	platform, err := platforms.Parse(architecture)
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}
	return platforms.Format(normalizeContainerPlatfrom(platform))
}

func GetUserContainerPlatform() string {
	return platforms.Format(getUserContainerPlatform())
}
