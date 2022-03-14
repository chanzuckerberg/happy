package cmd

import (
	"regexp"

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
	if notOk, err := stackNameIsInDnsCharset(args[0]); err != nil || notOk {
		return errors.New("STACK_NAME must only contain letters, digits, or hyphens, may not be all digits, and may not start or end with a hyphen")
	}
	return nil
}

func stackNameIsInDnsCharset(stackName string) (bool, error) {
	nonLdhPattern := "([^a-zA-Z0-9/-])"
	leadTrailHyphenPattern := "(^-|-$)"
	allDigitsPattern := "(^[0-9]*[0-9]$)"

	pattern := nonLdhPattern + "|" + leadTrailHyphenPattern + "|" + allDigitsPattern

	invalid, err := regexp.MatchString(pattern, stackName)
	return invalid, err
}
