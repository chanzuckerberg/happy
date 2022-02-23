package util

import (
	"context"
	"os/exec"
	"strings"

	semver "github.com/Masterminds/semver/v3"
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
	dockerComposeVersion := semver.MustParse(version)
	valid, reasons := dockerComposeMinVersion.Validate(dockerComposeVersion)
	if !valid {
		errs = multierror.Append(
			errs,
			errors.Errorf("docker compose >= V2 required but %s was detected", version),
		)
		for _, reason := range reasons {
			errs = multierror.Append(errs, reason)
		}
	}

	_, err = exec.LookPath("aws")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find aws cli in path"))
	}

	_, err = exec.LookPath("terraform")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find terraform in path"))
	}

	_, err = exec.LookPath("session-manager-plugin")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find session-manager-plugin in path"))
	}

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "docker engine is not running"))
	}

	_, err = client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "cannot connect to docker engine"))
	}

	return errs.ErrorOrNil()
}
