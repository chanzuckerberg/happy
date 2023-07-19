package githubconnector

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (client *GithubConnector) DownloadPackage(versionTag, os, arch, path string) (string, error) {

	release, err := client.GetRelease(versionTag)

	if err != nil {
		fmt.Println("Error getting release: ", err)
		return "", err
	}

	for _, asset := range release.Assets {
		if asset.Component == "happy" && asset.OS == os && asset.Architecture == arch {
			return download(asset.URL, asset.Name, path)
		}
	}

	return "", errors.New("no suitable package found")
}

func download(url, fileName, path string) (string, error) {

	// Download the url to the specified path
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
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
		return "", err
	}
	defer file.Close()

	return fileName, nil // TODO: return the full path to the downloaded file
}
