package util

import "fmt"

var (
	ReleaseVersion = "undefined"
	ReleaseGitSha  = "undefined"
)

type Release struct {
	Version string
	GitSha  string
}

func (r *Release) String() string {
	return fmt.Sprintf("version: %s, git_sha: %s", r.Version, r.GitSha)
}

func GetVersion() *Release {
	return &Release{
		Version: ReleaseVersion,
		GitSha:  ReleaseGitSha,
	}
}
