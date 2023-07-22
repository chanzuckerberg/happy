package githubconnector

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func (client *GithubConnector) DownloadPackage(version, os, arch, path string) (string, error) {

	release, err := client.GetRelease(version)

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

func download(url, fileName, dir string) (string, error) {

	filePath := path.Join(dir, fileName)
	// Download the url to the specified path/fileName
	file, err := os.Create(filePath)
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

	return filePath, nil
}
