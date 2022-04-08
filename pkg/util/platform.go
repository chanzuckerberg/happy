package util

import (
	"fmt"

	"github.com/containerd/containerd/platforms"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func normalizePlatfrom(platform v1.Platform) v1.Platform {
	platform.OS = "linux"
	return platforms.Normalize(platform)
}

func getUserPlatform() v1.Platform {
	platform := platforms.DefaultSpec()
	return normalizePlatfrom(platform)
}

func GetSystemPlatform(architecture string) string {
	platform, err := platforms.Parse(architecture)
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}
	return platforms.Format(normalizePlatfrom(platform))
}

func GetUserPlatform() string {
	return platforms.Format(getUserPlatform())
}
