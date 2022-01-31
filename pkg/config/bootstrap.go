package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	flagHappyProjectRoot        = "happy-project-root"
	flagHappyConfigPath         = "happy-config-path"
	flagDockerComposeConfigPath = "docker-compose-config-path"
)

// We will load bootrap configuration common to all commands here
// can then be consumed by other commands as needed.
var (
	happyProjectRoot        string
	happyConfigPath         string
	dockerComposeConfigPath string

	validate *validator.Validate
)

func init() {
	// use a single instance of Validate, it caches struct info
	validate = validator.New()
}

// RequireBootstrap wraps a command adding flags
// to resolve bootstrap configuration.
// NOTE that these can also be set by the environment
// and follow a pre-established convention of precedence.
// NOTE this should typically be called in a cobra commands init sequence.
func ConfigureCmdWithBootstrapConfig(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&happyProjectRoot, flagHappyProjectRoot, "", "Specify the root of your Happy project")
	cmd.PersistentFlags().StringVar(&happyConfigPath, flagHappyConfigPath, "", "Specify the path to your Happy project's config file")
	cmd.PersistentFlags().StringVar(&dockerComposeConfigPath, flagDockerComposeConfigPath, "", "Specify the path to your Happy project's docker compose file")
}

type Bootstrap struct {
	HappyConfigPath  string `envconfig:"HAPPY_CONFIG_PATH" validate:"required"`
	HappyProjectRoot string `envconfig:"HAPPY_PROJECT_ROOT" validate:"required"`

	DockerComposeConfigPath string `envconfig:"DOCKER_COMPOSE_CONFIG_PATH" validate:"required"`

	// TODO: do we want this overrideable?
	// TODO: Since it was hardcoded, leave as is and figure out later
	Env string `validate:"required,eq=rdev"`
}

func (b *Bootstrap) GetEnv() string {
	return b.Env
}

func (b *Bootstrap) GetHappyConfigPath() string {
	return b.HappyConfigPath
}

func (b *Bootstrap) GetHappyProjectRootPath() string {
	return b.HappyProjectRoot
}

func (b *Bootstrap) GetDockerComposeConfigPath() string {
	return b.DockerComposeConfigPath
}

func NewBootstrapConfig() (*Bootstrap, error) {
	// We compose this object going from lowest binding to strongest binding
	// overwriting as we go.
	// Once we've done all our steps, we will run a round of validation to make sure we have enough information

	// 1 - Default values
	b := &Bootstrap{
		// TODO(el): figure out why this is default and non-overwriteable
		Env: "rdev",
	}

	// 2 - environment variables
	err := envconfig.Process("", b)
	if err != nil {
		return nil, errors.Wrap(err, "could not read configuration from environment")
	}

	// 3 - CLI flags
	if happyProjectRoot != "" {
		b.HappyProjectRoot = happyProjectRoot
	}
	if happyConfigPath != "" {
		b.HappyConfigPath = happyConfigPath
	}
	if dockerComposeConfigPath != "" {
		b.DockerComposeConfigPath = dockerComposeConfigPath
	}

	// run validation
	err = validate.Struct(b)
	if err != nil {
		return nil, errors.Wrap(err, "invalid bootstrap configuration")
	}

	// expand paths to make it easier to consume
	b.HappyProjectRoot, err = homedir.Expand(b.HappyProjectRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not expand happy project root")
	}

	b.HappyConfigPath, err = homedir.Expand(b.HappyConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not expand happy config path")
	}

	b.DockerComposeConfigPath, err = homedir.Expand(b.DockerComposeConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not expand docker compose config path")
	}

	return b, nil
}
