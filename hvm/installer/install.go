package installer

import (
	"fmt"

	"github.com/chanzuckerberg/happy/shared/githubconnector"
)

func InstallPackage(versionTag, os, arch, binPath string) error {

	client := githubconnector.NewConnectorClient()

	downloaded, err := client.DownloadPackage(versionTag, os, arch, "/tmp")
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

	fmt.Println("Installing package from ", sourcePackagePath, " to ", binPath)

	return nil
}
