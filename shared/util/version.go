package util

import (
	"fmt"
)

var (
	ReleaseVersion = "undefined"
	ReleaseGitSha  = "undefined"
)

type Release struct {
	Version string
	GitSha  string
}

func (r *Release) String() string {
	return fmt.Sprintf("version: %s\ngit_sha: %s", r.Version, r.GitSha)
}

func (r *Release) Equal(otherRelease *Release) bool {
	return r.GitSha == otherRelease.GitSha
}

func GetVersion() *Release {
	return &Release{
		Version: ReleaseVersion,
		GitSha:  ReleaseGitSha,
	}
}
