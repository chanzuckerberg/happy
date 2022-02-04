package config

import (
	"context"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CommandValidator func(ctx context.Context) error

func RequireDocker(ctx context.Context) error {
	required := []string{"docker-compose", "docker"}
	for _, r := range required {
		_, err := exec.LookPath(r)
		if err != nil {

			return errors.Wrapf(err, "could not find %s executable in PATH", r)
		}
	}
	return nil
}

func WithPreRunValidation(validators ...CommandValidator) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		for _, validator := range validators {
			err := validator(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}

}
