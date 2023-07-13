package github_connector

import (
	"github.com/google/go-github/v53/github"
)

type GithubConnector struct {
	github *github.Client
}

type Release struct {
	Tag     string
	Version string
	Assets  []ReleaseAsset
}

type ReleaseAsset struct {
	Name         string
	OS           string
	Architecture string
	URL          string
	FileType     string
}

func NewConnectorClient() *GithubConnector {
	return &GithubConnector{
		github: github.NewClient(nil),
	}
}
