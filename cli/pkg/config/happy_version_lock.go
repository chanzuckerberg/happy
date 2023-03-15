package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type HappyVersionLockFile struct {
	HappyVersion string
	path         string
}

func NewHappyVersionLockFile(projectRoot string) *HappyVersionLockFile {
	path := calcHappyVersionPath(projectRoot)

	return &HappyVersionLockFile{
		HappyVersion: "",
		path:         path,
	}
}

func LoadHappyVersionLockFile(projectRoot string) (*HappyVersionLockFile, error) {
	return NewHappyVersionLockFile(projectRoot), nil
}

func (v *HappyVersionLockFile) Save() error {

	path, err := v.GetPath()
	if err != nil {
		return err
	}

	happyVersionFile, err := os.Create(path)

	if err != nil {
		return errors.New(fmt.Sprintf("Could not create %s: %v", v.path, err))
	}

	happyVersionFile.WriteString(v.HappyVersion)
	happyVersionFile.Close()

	return nil
}

func calcHappyVersionPath(projectRoot string) string {
	versionFilePath := filepath.Join(projectRoot, ".happy", "version.lock")
	return versionFilePath
}

func (v *HappyVersionLockFile) SetVersion(version string) error {

	if version == "" {
		return errors.New("Empty version is not allowed")
	}

	v.HappyVersion = version

	return nil
}

func (v *HappyVersionLockFile) GetVersion() (string, error) {

	if v.HappyVersion == "" {
		return "", errors.New("Version is not set")
	}

	return v.HappyVersion, nil
}

func (v *HappyVersionLockFile) GetPath() (string, error) {
	if v.path == "" {
		return "", errors.New("Path is not set")
	}
	return v.path, nil
}
