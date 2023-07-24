package linkmanager

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func SetBinLink(org, project, version string) error {

	user, err := user.Current()

	if err != nil {
		return errors.Wrapf(err, "Error getting current user information")
	}

	home := user.HomeDir
	versionsPath := path.Join(home, ".czi", "versions", org, project, version)
	binPath := path.Join(home, ".czi", "bin")

	os.MkdirAll(binPath, 0755)

	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		fmt.Println("Requested version is not installed.")
		return err
	}

	// Iterate through all the files in versionsPath

	files, err := os.ReadDir(versionsPath)
	if err != nil {
		return errors.Wrapf(err, "Error reading directory")
	}

	var bin string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return errors.Wrapf(err, "Error getting file info for %s", file.Name())
		}

		fmt.Println("Checking ", file.Name(), " permissions", info.Mode())
		// Skip if the file is not owner-executable
		if !strings.Contains(info.Mode().String(), "x") {
			fmt.Println("Skipping ", file.Name(), " as it is not executable")
			continue
		}

		bin = file.Name()

		fmt.Println("Setting bin link for ", bin)

		os.Remove(path.Join(binPath, bin))

		err = os.Symlink(path.Join(versionsPath, bin), path.Join(binPath, bin))
		if err != nil {
			return errors.Wrapf(err, "Error creating symlink")
		}

	}

	return nil
}
