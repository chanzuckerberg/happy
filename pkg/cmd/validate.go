package cmd

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
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

func CheckStackName(cmd *cobra.Command, args []string) error {
	// return anonymous function parameterized on arg position instead? :thinking:
	if stackNameIsInDnsCharset(args[0]) {
		return nil
	} else {
		return errors.New("STACK_NAME must only contain letters, digits, or hyphens")
	}
}

func stackNameIsInDnsCharset(stackName string) bool {
	// TODO
	return true
}
