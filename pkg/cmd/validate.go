package cmd

import (
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

type Validate func(cmd *cobra.Command, args []string) error

func ValidateMany(vs ...Validate) Validate {
	return func(cmd *cobra.Command, args []string) error {
		var err *multierror.Error
		for _, v := range vs {
			err = multierror.Append(err, v(cmd, args))
		}
		return err.ErrorOrNil()
	}
}
