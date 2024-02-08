package util

import (
	"path/filepath"

	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func ensureDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func GetCachePath() (string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "unable to find home directory")
	}
	dir := filepath.Join(homedir, ".happy", "cache")
	if err := ensureDir(dir); err != nil {
		return "", errors.Wrap(err, "unable to create cache directory")
	}
	return dir, nil
}
