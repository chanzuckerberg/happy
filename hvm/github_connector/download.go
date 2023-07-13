package github_connector

import "fmt"

func (client *GithubConnector) DownloadPackage(versionTag, os, arch, path string) error {
	fmt.Printf("Downloading %s for %s/%s to %s\n", versionTag, os, arch, path)
	return nil
}
