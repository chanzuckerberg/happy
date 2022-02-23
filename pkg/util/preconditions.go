package util

import (
	"context"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func ValidateEnvironment(ctx context.Context) error {
	var errs error
	_, err := exec.LookPath("docker")
	if err != nil {
		errs = multierror.Append(errs, errors.Wrap(err, "could not find docker in path"))
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

	return errs
}
