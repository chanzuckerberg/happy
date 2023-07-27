package installer

import (
	"bytes"
	"context"
	"io/fs"
	"os"

	"github.com/chanzuckerberg/go-misc/errors"
	"github.com/chanzuckerberg/happy/shared/githubconnector"
	"github.com/codeclysm/extract"
)

func InstallPackage(ctx context.Context, org, project, version, opsys, arch, binPath string) error {

	client := githubconnector.NewConnectorClient()

	downloaded, err := client.DownloadPackage(org, project, version, opsys, arch, "/tmp")
	if err != nil {
		return errors.Wrap(err, "downloading package")
	}

	err = doInstall(ctx, downloaded, binPath)

	if err != nil {
		return errors.Wrap(err, "installing package")
	}

	return nil
}

func doInstall(ctx context.Context, sourcePackagePath, binPath string) error {

	err := os.MkdirAll(binPath, fs.FileMode(0755))
	if err != nil {
		return errors.Wrapf(err, "Error creating directory %s", binPath)
	}

	data, _ := os.ReadFile(sourcePackagePath)
	buffer := bytes.NewBuffer(data)
	err = extract.Gz(ctx, buffer, binPath, nil)
	os.Remove(sourcePackagePath)

	if err != nil {
		return errors.Wrapf(err, "extracting package %s", sourcePackagePath)
	}

	return nil
}
