package linkmanager

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func SetBinLink(org, project, version string) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrapf(err, "Error getting home directory")
	}

	versionsPath := path.Join(home, ".czi", "versions", org, project, version)
	binPath := path.Join(home, ".czi", "bin")

	err = os.MkdirAll(binPath, 0755)
	if err != nil {
		return errors.Wrapf(err, "creating directory %s", binPath)
	}

	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		return errors.Wrap(err, "requested version is not installed")
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
			logrus.Printf("Skipping %s as it is not executable", file.Name())
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
