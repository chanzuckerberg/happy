package githubconnector

import (
	"github.com/google/go-github/v53/github"
    "golang.org/x/oauth2"
    "context"
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

func NewConnectorClient(githubPAT *string) *GithubConnector {

    if githubPAT == nil {
        return &GithubConnector{
            github: github.NewClient(nil),
        }
    }

    ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubPAT},
	)
	tc := oauth2.NewClient(ctx, ts)

    client := github.NewClient(tc)

	return &GithubConnector{
        github: client,
    }
}
