package config

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagHappyProjectRoot = "project-root"
	flagHappyConfigPath  = "config-path"
	flagEnv              = "env"
	FlagAWSProfile       = "aws-profile"

	flagComposeEnvFile          = "docker-compose-env-file"
	flagDockerComposeConfigPath = "docker-compose-config-path"
)

// We will load bootrap configuration common to all commands here
// can then be consumed by other commands as needed.
var (
	happyProjectRoot        string
	happyConfigPath         string
	dockerComposeConfigPath string
	env                     string
	composeEnvFile          string
	awsProfile              string

	errCouldNotInferFindHappyRoot = errors.New("could not infer .happy root")

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
	cmd.PersistentFlags().StringVar(&env, flagEnv, "", "Specify a Happy env")
	cmd.PersistentFlags().StringVar(&composeEnvFile, flagComposeEnvFile, "", "Environment file to pass to docker compose")
	cmd.PersistentFlags().StringVar(&awsProfile, FlagAWSProfile, "", "Override the AWS profile to use. If speficied but empty, will use the default credentil chain.")
}

type Bootstrap struct {
	HappyConfigPath  string `envconfig:"HAPPY_CONFIG_PATH" validate:"required"`
	HappyProjectRoot string `envconfig:"HAPPY_PROJECT_ROOT" validate:"required"`

	DockerComposeConfigPath  string `envconfig:"DOCKER_COMPOSE_CONFIG_PATH" validate:"required"`
	DockerComposeEnvFilePath string `envconfig:"DOCKER_COMPOSE_ENV_FILE_PATH"`

	AWSProfile *string `envconfig:"AWS_PROFILE"`
	AWSRegion  *string `envconfig:"AWS_REGION"`
	AWSRoleARN *string

	Env string `envconfig:"HAPPY_ENV"`
}

func (b *Bootstrap) GetEnv() string {
	return b.Env
}

func (b *Bootstrap) GetComposeEnvFile() string {
	return composeEnvFile
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

func (b *Bootstrap) GetAWSProfile() *string {
	return b.AWSProfile
}

func (b *Bootstrap) GetAWSRegion() *string {
	return b.AWSRegion
}

// We search up the directory structure until we find we are
// in a directory that contains a .happy dir
func searchHappyRoot(path string) (string, error) {
	if path == "/" {
		return "", errCouldNotInferFindHappyRoot
	}

	potentialHappyDir := filepath.Join(path, "/.happy")
	logrus.Debugf("searching happy root at %s", potentialHappyDir)
	_, err := os.Stat(potentialHappyDir)
	// if not here, keep going up
	if errors.Is(err, fs.ErrNotExist) {
		dir := filepath.Dir(path)
		return searchHappyRoot(dir)
	}
	// If we get a permission denied, stop
	if errors.Is(err, fs.ErrPermission) {
		return "", errors.Wrap(err, errCouldNotInferFindHappyRoot.Error())
	}
	// other errors, bubble them up
	if err != nil {
		return "", errors.Wrap(err, "unexpected err while searching for .happy root")
	}
	// if no error, we found what we're looking for
	return path, nil
}

func NewBootstrapConfig(cmd *cobra.Command) (*Bootstrap, error) {
	return newBootstrap(env, cmd.Flags().Changed(FlagAWSProfile))
}

func NewBootstrapConfigForEnv(env string, useAWSProfile bool) (*Bootstrap, error) {
	return newBootstrap(env, useAWSProfile)
}

// This is a simple bootstrap used for the bootstrapping command
func NewSimpleBootstrap(cmd *cobra.Command) (*Bootstrap, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "could not get working directory")
	}
	b := &Bootstrap{
		Env:              "",
		HappyProjectRoot: path,
	}

	err = envconfig.Process("", b)
	if err != nil {
		return nil, errors.Wrap(err, "could not process env vars")
	}

	if b.HappyConfigPath == "" {
		b.HappyConfigPath = filepath.Join(b.HappyProjectRoot, "/.happy/config.json")
	}

	if b.DockerComposeConfigPath == "" {
		b.DockerComposeConfigPath = filepath.Join(b.HappyProjectRoot, "/docker-compose.yml")
	}

	return b, nil
}

func newBootstrap(env string, useAWSProfile bool) (*Bootstrap, error) {
	// We compose this object going from the lowest binding to the strongest binding
	// overwriting as we go.
	// Once we've done all our steps, we will run a round of validation to make sure we have enough information

	// 0 - preamble, gather background info
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "could not get working directory")
	}

	defaultHappyRoot, err := searchHappyRoot(wd)
	if err != nil && !errors.Is(err, errCouldNotInferFindHappyRoot) {
		return nil, err
	}

	// 1 - Default values
	b := &Bootstrap{
		Env:              "",
		HappyProjectRoot: defaultHappyRoot,
	}

	// 2 - environment variables
	err = envconfig.Process("", b)
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
	// Bootstrap Env will be read from envconfig if it exists
	// If the --env flag was set, the flag will override that value
	// Otherwise it will pass through as an empty string and depend on happy config default_env
	if env != "" {
		b.Env = env
	}

	// NOTE: We treat "" profile as asking to use the default provider chain
	if useAWSProfile {
		b.AWSProfile = &awsProfile
	}

	// 4 - Inferred
	// These are like defaults but rely on info we've gathered so far
	if b.HappyConfigPath == "" {
		b.HappyConfigPath = filepath.Join(b.HappyProjectRoot, "/.happy/config.json")
	}

	if b.DockerComposeConfigPath == "" {
		b.DockerComposeConfigPath = filepath.Join(b.HappyProjectRoot, "/docker-compose.yml")
	}

	// run validation
	err = validate.Struct(b)
	if err != nil {
		return nil, prettyValidationErrors(err)
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

	// validate paths exist
	_, err = os.Stat(b.DockerComposeConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "docker compose config not found at %s", b.DockerComposeConfigPath)
	}
	_, err = os.Stat(b.HappyConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "happy config not found at %s", b.HappyConfigPath)
	}

	return b, nil
}

func prettyValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	var originalErrs validator.ValidationErrors
	if !errors.As(err, &originalErrs) {
		return err
	}
	var errs error
	for _, err := range originalErrs {
		niceErr := errors.Errorf("%s is required but was not set and could not be inferred", err.Field())
		errs = multierror.Append(errs, niceErr)
	}
	return errs
}
