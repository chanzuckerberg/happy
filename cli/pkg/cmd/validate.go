package cmd

import (
	"regexp"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Validate(vs ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, v := range vs {
			err := v(cmd, args)
			if err != nil {
				return multierror.Append(err, cmd.Help())
			}
		}
		return nil
	}
}

func IsImageEnvUsedWithImageStack(cmd *cobra.Command, args []string) error {
	imageSrcEnv, err := cmd.Flags().GetString(flagImageSrcEnv)
	if err != nil {
		return err
	}
	imageSrcStack, err := cmd.Flags().GetString(flagImageSrcStack)
	if err != nil {
		return err
	}
	if imageSrcEnv == "" && imageSrcStack == "" {
		return nil
	}

	if cmd.Flags().Changed(flagImageSrcEnv) && !cmd.Flags().Changed(flagImageSrcStack) {
		return errors.New("--image-src-env can only be used when --image-src-stack is specified")
	}
	if !cmd.Flags().Changed(flagImageSrcEnv) && cmd.Flags().Changed(flagImageSrcStack) {
		return errors.New("--image-src-env can only be used when --image-src-stack is specified")
	}

	return nil
}

func IsTagUsedWithSkipTag(cmd *cobra.Command, args []string) error {
	createTag, err := cmd.Flags().GetBool("create-tag")
	if err != nil {
		return err
	}
	if cmd.Flags().Changed("skip-check-tag") && !cmd.Flags().Changed("tag") {
		return errors.New("--skip-check-tag can only be used when --tag is specified")
	}

	if !createTag && !cmd.Flags().Changed("tag") {
		return errors.New("Must specify a tag when create-tag=false")
	}

	return nil
}

func IsStackNameDNSCharset(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New(("Command does not contain a STACK_NAME"))
	}
	if notOk, err := stackNameIsInDnsCharset(args[0]); err != nil || notOk {
		return errors.New("STACK_NAME must only contain letters, digits, or hyphens, may not be all digits, and may not start or end with a hyphen")
	}
	return nil
}

func IsStackNameAlphaNumeric(cmd *cobra.Command, args []string) error {
	if !regexp.MustCompile(`^[a-z0-9\-]*$`).MatchString(args[0]) {
		return errors.New("stack name must contain only lowercase letters, numbers, and hyphens/dashes")
	}
	return nil
}

func stackNameIsInDnsCharset(stackName string) (bool, error) {
	nonLdhPattern := "([^a-zA-Z0-9/-])"
	leadTrailHyphenPattern := "(^-|-$)"
	allDigitsPattern := "(^[0-9]*[0-9]$)"
	overLengthPattern := ".{64,}"

	pattern := nonLdhPattern + "|" + leadTrailHyphenPattern + "|" + allDigitsPattern + "|" + overLengthPattern

	invalid, err := regexp.MatchString(pattern, stackName)
	return invalid, err
}
