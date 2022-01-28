package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

type Bootstrap struct {
	happyConfigPath  string `envconfig:"HAPPY_CONFIG_PATH"`
	happyProjectRoot string `envconfig:"HAPPY_PROJECT_ROOT"`

	dockerComposeConfigPath string `envconfig:"DOCKER_COMPOSE_CONFIG_PATH"`

	// TODO: do we want this overrideable? For now it was hardcoded and not something we can change so leaving as is
	env string
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
	// We compose this object going from lowest binding to strongest binding
	// overwriting as we go.
	// Once we've done all our steps, we will run a round of validation to make sure we have enough information

	// 1 - Default values
	b := &Bootstrap{
		// TODO(el): figure out why this is default and non-overwriteable
		env: "rdev",
	}

	// 2 - environment variables

	return nil, nil
}
