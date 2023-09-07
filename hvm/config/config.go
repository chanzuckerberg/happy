package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"path"
)

type HvmConfig struct {
	GithubPAT *string
}

func GetHvmConfig() (*HvmConfig, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return nil, errors.Wrap(err, "getting current user home directory")
	}

	configPath := path.Join(home, ".czi", "etc", "hvmconfig.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "loading config file")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "opening config file")
	}
	defer file.Close()

	// Parse json from file into HvmConfig struct

	output := &HvmConfig{}
	err = json.NewDecoder(file).Decode(&output)

	if err != nil {
		return nil, errors.Wrap(err, "parsing config file")
	}

	// Return HvmConfig struct

	return output, nil

}
