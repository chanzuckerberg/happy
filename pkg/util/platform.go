package util

import (
	"fmt"
	"runtime"

	"github.com/containerd/containerd/platforms"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func getUserPlatform() v1.Platform {
	platform := platforms.DefaultSpec()
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		platform.OS = "linux"
	}
	return platforms.Normalize(platform)
}

func GetSystemPlatform(architecture string) string {
	platform, err := platforms.Parse(architecture)
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}
	platform.OS = "linux"
	return platforms.Format(platforms.Normalize(platform))
}

func GetUserPlatform() string {
	return platforms.Format(getUserPlatform())
}
