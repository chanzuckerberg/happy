package githubconnector

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"io"
	"net/http"
	"os"
	"path"
)

func (client *GithubConnector) DownloadPackage(org, project, version, os, arch, path string) (string, error) {

	release, err := client.GetRelease(org, project, version)

	if err != nil {
		return "", errors.Wrap(err, "getting release")
	}

	for _, asset := range release.Assets {
		if asset.Component == project && asset.OS == os && asset.Architecture == arch {
			return download(asset.URL, asset.Name, path)
		}
	}

	return "", errors.New("no suitable package found")
}

func download(url, fileName, dir string) (string, error) {

	filePath := path.Join(dir, fileName)
	// Download the url to the specified path/fileName

	logrus.Debugf("Downloading %s to %s\n", url, filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.Wrap(err, "creating download file")
	}
	defer file.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "fetching package")
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		os.Remove(fileName)
		return "", errors.Wrap(err, "writing download file")
	}

	return filePath, nil
}
