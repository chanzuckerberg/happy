package util

import (
	"context"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func ValidateEnvironment() error {
	_, err := exec.LookPath("docker")
	if err != nil {
		return errors.Wrap(err, "could not find docker in path")
	}

	_, err = exec.LookPath("aws")
	if err != nil {
		return errors.Wrap(err, "could not find aws cli in path")
	}

	_, err = exec.LookPath("terraform")
	if err != nil {
		return errors.Wrap(err, "could not find terraform in path")
	}

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "docker engine is not running")
	}

	_, err = client.ContainerList(context.Background(), types.ContainerListOptions{})
	return errors.Wrap(err, "cannot connect to docker engine")
}
