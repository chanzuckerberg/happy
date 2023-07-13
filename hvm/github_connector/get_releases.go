package github_connector

import (
	"context"
	"strings"

	"github.com/google/go-github/v53/github"
)

func (client *GithubConnector) GetRelease(versionTag string) (*Release, error) {

	ghRelease, _, err := client.github.Repositories.GetReleaseByTag(context.Background(), "chanzuckerberg", "happy", "v0.0.1")

	if err != nil {
		return nil, err
	}

	return &Release{
		Tag:     *ghRelease.TagName,
		Version: strings.Replace(*ghRelease.TagName, "v", "", 1),
	}, nil

}

func (client *GithubConnector) GetReleases(org string, repo string) ([]*Release, error) {

	happyReleases := make([]*Release, 0)

	// Only up to 1000 results are returned by the API
	for page := 1; page < 2; page++ {
		releases, _, err := client.github.Repositories.ListReleases(context.TODO(), org, repo,
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

			assets := make([]ReleaseAsset, 0)

			for _, asset := range release.Assets {

				assets = append(assets, ReleaseAsset{
					Name:         asset.GetName(),
					OS:           nameToOS(asset.GetName()),
					Architecture: nameToArchitecture(asset.GetName()),
					URL:          asset.GetBrowserDownloadURL(),
					FileType:     asset.GetContentType(),
				})

			}

			if strings.HasPrefix(*release.TagName, "v") {
				happyReleases = append(happyReleases, &Release{
					Tag:     *release.TagName,
					Version: tagToVersion(*release.TagName),
					Assets:  assets,
				})
			}
		}
	}

	return happyReleases, nil

}

func tagToVersion(tag string) string {
	return strings.Replace(tag, "v", "", 1)
}

func nameToArchitecture(name string) string {
	parts := strings.Split(name, "_")

	if len(parts) < 4 {
		return ""
	}

	archAndExtension := strings.Split(name, "_")[3]
	return strings.Split(archAndExtension, ".")[0]
}

func nameToOS(label string) string {
	os := strings.Split(label, "_")[2]

	if os != "checksums.txt" {
		return os
	}

	return ""
}
