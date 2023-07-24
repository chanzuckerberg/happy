package installer

import (
	"bytes"
	"context"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/chanzuckerberg/go-misc/errors"
	"github.com/chanzuckerberg/happy/shared/githubconnector"
	"github.com/codeclysm/extract"
)

func InstallPackage(org, project, version, os, arch, binPath string) error {

	client := githubconnector.NewConnectorClient()

	downloaded, err := client.DownloadPackage(org, project, version, os, arch, "/tmp")
	if err != nil {
		return err
	}

	err = doInstall(downloaded, binPath)

	if err != nil {
		return err
	}

	return nil
}

func doInstall(sourcePackagePath, binPath string) error {

	err := os.MkdirAll(binPath, fs.FileMode(0755))
	if err != nil {
		return errors.Wrapf(err, "Error creating directory %s", binPath)
	}

	data, _ := ioutil.ReadFile(sourcePackagePath)
	buffer := bytes.NewBuffer(data)
	err = extract.Gz(context.TODO(), buffer, binPath, nil)

	if err != nil {
		return errors.Wrapf(err, "Error extracting package %s", sourcePackagePath)
	}

	return nil
}
