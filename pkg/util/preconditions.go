package util

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func ValidateEnvironment(ctx context.Context) error {
	dockerComposeMinVersion, err := semver.NewConstraint(">= v2")
	if err != nil {
		return errors.Wrap(err, "could not establish docker compose version")
	}

	var errs *multierror.Error
	_, err = exec.LookPath("docker")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find docker in path"))
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

	_, err = exec.LookPath("aws")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find aws cli in path, run 'brew install awscli' or follow https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html"))
	}

	_, err = exec.LookPath("terraform")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find terraform cli in path, run 'brew install terraform', or follow https://learn.hashicorp.com/tutorials/terraform/install-cli"))
	}

	_, err = exec.LookPath("session-manager-plugin")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find session-manager-plugin in path, run 'brew install --cask session-manager-plugin', or follow https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html"))
	}

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "docker engine is not running, follow https://docs.docker.com/get-docker/"))
	}

	_, err = client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "cannot connect to docker engine"))
	}

	return errs.ErrorOrNil()
}
