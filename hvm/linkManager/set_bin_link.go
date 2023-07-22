package linkmanager

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/pkg/errors"
)

func SetBinLink(version string) error {

	user, err := user.Current()

	if err != nil {
		return errors.Wrapf(err, "Error getting current user information")
	}

	home := user.HomeDir
	versionsPath := path.Join(home, ".czi", "versions", "happy", version)
	binPath := path.Join(home, ".czi", "bin")

	os.MkdirAll(binPath, 0755)

	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		fmt.Println("Requested version is not installed.")
		return err
	}

	os.Remove(path.Join(binPath, "happy"))

	err = os.Symlink(path.Join(versionsPath, "happy"), path.Join(binPath, "happy"))
	if err != nil {
		return errors.Wrapf(err, "Error creating symlink")
	}

	return nil
}
