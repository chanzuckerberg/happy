package github_connector

import (
	"context"
	"strings"

	"github.com/google/go-github/v53/github"
)

func (client *GithubConnector) GetHappyRelease(versionTag string) (*Release, error) {

	ghRelease, _, err := client.github.Repositories.GetReleaseByTag(context.Background(), "chanzuckerberg", "happy", "v0.0.1")

	if err != nil {
		return nil, err
	}

	return &Release{
		Tag:     *ghRelease.TagName,
		Version: strings.Replace(*ghRelease.TagName, "v", "", 1),
	}, nil

}

func (client *GithubConnector) GetHappyReleases() ([]*Release, error) {

	happyReleases := make([]*Release, 0)

	// Only up to 1000 results are returned by the API
	for page := 1; page < 10; page++ {
		releases, _, err := client.github.Repositories.ListReleases(context.TODO(), "chanzuckerberg", "happy",
			&github.ListOptions{
				Page:    page,
				PerPage: 100,
			})

		if err != nil {
			return nil, err
		}

		if len(releases) == 0 {
			break
		}

		for _, release := range releases {

			if strings.HasPrefix(*release.TagName, "v") {
				happyReleases = append(happyReleases, &Release{
					Tag:     *release.TagName,
					Version: tagToVersion(*release.TagName),
				})
			}
		}
	}

	return happyReleases, nil

}

func tagToVersion(tag string) string {
	return strings.Replace(tag, "v", "", 1)
}
