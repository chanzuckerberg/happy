package cmd

import (
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

func Validate(vs ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		var err *multierror.Error
		for _, v := range vs {
			err = multierror.Append(err, v(cmd, args))
		}
		// If argument validation fails, print help
		if err.ErrorOrNil() != nil {
			return multierror.Append(err, cmd.Help())
		}
		return nil
	}
}
