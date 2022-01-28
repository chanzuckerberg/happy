package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

type Bootstrap struct {
	happyConfigPath  string
	happyProjectRoot string

	env string

	dockerComposeConfigPath string
}

func (b *Bootstrap) GetEnv() string {
	return b.env
}

func (b *Bootstrap) GetHappyConfigPath() (string, error) {
	expanded, err := homedir.Expand(b.happyConfigPath)
	return expanded, errors.Wrap(err, "could not expand happy config path")
}

func (b *Bootstrap) GetHappyProjectRootPath() (string, error) {
	expanded, err := homedir.Expand(b.happyProjectRoot)
	return expanded, errors.Wrap(err, "could not expand happy root path")
}

func (b *Bootstrap) GetDockerComposeConfigPath() (string, error) {
	expanded, err := homedir.Expand(b.dockerComposeConfigPath)
	return expanded, errors.Wrap(err, "could not expand docker compose config path")
}

func ResolveBootstrapConfig() (*Bootstrap, error) {
	return nil, nil
}
