package githubconnector

import (
	"github.com/google/go-github/v53/github"
    "golang.org/x/oauth2"
    "context"
    "github.com/chanzuckerberg/happy/hvm/config"
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
	Component    string
	OS           string
	Architecture string
	URL          string
	FileType     string
}

func NewConnectorClient() *GithubConnector {

    hvmConfig, _ := config.GetHvmConfig()

    // If we don't get a config, just don't use a PAT.
    // No need to bail out.

    if hvmConfig == nil || hvmConfig.GithubPAT == nil {
        return &GithubConnector{
            github: github.NewClient(nil),
        }
    }

    ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *hvmConfig.GithubPAT},
	)
	tc := oauth2.NewClient(ctx, ts)

    client := github.NewClient(tc)

	return &GithubConnector{
        github: client,
    }
}
