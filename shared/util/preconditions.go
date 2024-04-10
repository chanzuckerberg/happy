package util

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ValidationCallback = func(context.Context) error

type ValidationCheckList struct {
	MinDockerComposeVersion          ValidationCallback
	DockerInstalled                  ValidationCallback
	DockerEngineRunning              ValidationCallback
	AwsInstalled                     ValidationCallback
	TerraformInstalled               ValidationCallback
	AwsSessionManagerPluginInstalled ValidationCallback
}

// Takes a list of callbacks to run. Canned ones are in the ValidationCheckList struct,
// but nothing stops you from adding custom ones from other packages.
func ValidateEnvironment(ctx context.Context, validations ...ValidationCallback) error {
	var errs *multierror.Error

	for _, validation := range validations {
		err := validation(ctx)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}

// Defaults run all validations
func NewValidationCheckList() *ValidationCheckList {
	return &ValidationCheckList{
		MinDockerComposeVersion:          ValidateMinDockerComposeVersion,
		DockerInstalled:                  ValidateDockerInstalled,
		DockerEngineRunning:              ValidateDockerEngineRunning,
		AwsInstalled:                     ValidateAwsInstalled,
		TerraformInstalled:               ValidateTerraformInstalled,
		AwsSessionManagerPluginInstalled: ValidateAwsSessionManagerPluginInstalled,
	}
}

func ValidateDockerInstalled(ctx context.Context) error {
	var errs *multierror.Error

	_, err := exec.LookPath("docker")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find docker in path"))
	}
	return errs.ErrorOrNil()
}

func ValidateMinDockerComposeVersion(ctx context.Context) error {
	var errs *multierror.Error

	dockerComposeMinVersion, err := semver.NewConstraint(">= v2.0.0-0")
	if err != nil {
		return errors.Wrap(err, "could not establish docker compose version")
	}

	v, err := exec.CommandContext(ctx, "docker", "compose", "version", "--short").Output()
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not determine docker compose version"))
	}

	version := strings.TrimSpace(string(v))
	dockerComposeVersion, err := semver.NewVersion(version)
	if err != nil {
		return errors.Wrapf(err, `invalid docker compose version. docker compose >= V2 required but "%s" was detected, please follow https://docs.docker.com/compose/cli-command/`, version)
	}
	valid, reasons := dockerComposeMinVersion.Validate(dockerComposeVersion)
	if !valid {
		errs = multierror.Append(
			errs,
			errors.Errorf("docker compose >= V2 required but %s was detected, please follow https://docs.docker.com/compose/cli-command/", version),
		)
		for _, reason := range reasons {
			errs = multierror.Append(errs, reason)
		}
	}

	return errs.ErrorOrNil()
}

func ValidateAwsInstalled(ctx context.Context) error {
	var errs *multierror.Error
	_, err := exec.LookPath("aws")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find aws cli in path, run 'brew install awscli' or follow https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html"))
	}

	return errs.ErrorOrNil()
}

func ValidateAwsSessionManagerPluginInstalled(ctx context.Context) error {
	var errs *multierror.Error
	_, err := exec.LookPath("session-manager-plugin")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find session-manager-plugin in path, run 'brew install --cask session-manager-plugin', or follow https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html"))
	}

	return errs.ErrorOrNil()
}

func ValidateTerraformInstalled(ctx context.Context) error {
	var errs *multierror.Error
	_, err := exec.LookPath("terraform")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find terraform cli in path, run 'brew install terraform', or follow https://learn.hashicorp.com/tutorials/terraform/install-cli"))
	}

	return errs.ErrorOrNil()
}

func ValidateDockerEngineRunning(ctx context.Context) error {
	var errs *multierror.Error

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "docker engine is not running, follow https://docs.docker.com/get-docker/"))
	}

	log.Debug("checking docker engine is running; if the process freezes up, please restart docker engine")
	_, err = client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "cannot connect to docker engine"))
	}

	return errs.ErrorOrNil()
}
