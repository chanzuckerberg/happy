package github_connector

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (client *GithubConnector) DownloadPackage(versionTag, os, arch, path string) error {
	fmt.Printf("Downloading %s for %s/%s to %s\n", versionTag, os, arch, path)

	release, err := client.GetRelease(versionTag)

	if err != nil {
		fmt.Println("Error getting release: ", err)
		return err
	}

	for _, asset := range release.Assets {
		if asset.Component == "happy" && asset.OS == os && asset.Architecture == arch {
			download(asset.URL, asset.Name, path)
			break
		}
	}

	return nil
}

func download(url, fileName, path string) error {
	fmt.Printf("Downloading %s to %s\n", url, path)

	// Download the url to the specified path
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		os.Remove(fileName)
		return err
	}
	defer file.Close()

	return nil
}
