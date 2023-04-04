package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type HappyVersionLockFile struct {
	HappyVersion string
	Path         string `json:"-"`
}

func NewHappyVersionLockFile(projectRoot string, requiredVersion string) (*HappyVersionLockFile, error) {
	return &HappyVersionLockFile{
		HappyVersion: requiredVersion,
		Path:         calcHappyVersionPath(projectRoot),
	}, nil
}

func DoesHappyVersionLockFileExist(projectRoot string) bool {
	filePath := calcHappyVersionPath(projectRoot)
	_, err := os.Stat(filePath)
	return err == nil
}

func LoadHappyVersionLockFile(projectRoot string) (*HappyVersionLockFile, error) {
	versionFile, err := os.Open(calcHappyVersionPath(projectRoot))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open happy version lock file")
	}
	defer versionFile.Close()

	hvlf := HappyVersionLockFile{}

	err = json.NewDecoder(versionFile).Decode(&hvlf)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing happy version lock file")
	}

	return &hvlf, nil
}

func (v *HappyVersionLockFile) Save() (err error) {
	contents, err := json.MarshalIndent(&v, "", " ")
	if err != nil {
		return errors.Wrap(err, "could not marshal config file contents")
	}

	happyVersionFile, err := os.Create(v.Path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not create %s", v.Path))
	}

	defer func() { err = happyVersionFile.Close() }()

	_, err = happyVersionFile.WriteString(string(contents))
	if err != nil {
		return errors.Wrap(err, "error writing to happy config version lock")
	}

	err = happyVersionFile.Sync()
	if err != nil {
		return errors.Wrap(err, "error syncing happy config version lock")
	}

	return nil
}

func calcHappyVersionPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".happy", "version.lock")
}
