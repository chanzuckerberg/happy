package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func SupportUpdateSlices(cmd *cobra.Command, sliceName *string, sliceDefaultTag *string) {
	SupportBuildSlices(cmd, sliceName, sliceDefaultTag)
	cmd.Flags().StringVar(sliceDefaultTag, "slice-default-tag", "", "For stacks using slices, override the default tag for any images that aren't being built & pushed by the slice")
}

func SupportBuildSlices(cmd *cobra.Command, sliceName *string, sliceDefaultTag *string) {
	cmd.Flags().StringVarP(sliceName, "slice", "s", "", "If you only need to test a slice of the app, specify it here")
}

func ValidateUpdateSliceFlags(cmd *cobra.Command, args []string) error {
	if !(cmd.Flags().Changed("slice") == cmd.Flags().Changed("slice-default-tag")) {
		return errors.New("both (or None) of `slice` and `slice-default-tag` must be set (or unset).")
	}
	return nil
}

const (
	flagDoRunMigrations = "do-migrations"
	flagSkipMigrations  = "skip-migrations"
)

func SetMigrationFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("do-migrations", true, "Specify if you want to force migrations to run")
	cmd.Flags().Bool("skip-migrations", false, "Specify if you want to skip migrations")
}

func ShouldRunMigrations(ctx context.Context, cmd *cobra.Command, happyConf *config.HappyConfig) (bool, error) {
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		return false, nil
	}
	if cmd.Flags().Changed(flagDoRunMigrations) && cmd.Flags().Changed(flagSkipMigrations) {
		return false, errors.Errorf(
			"flags %s and %s cannot be specified at the same time",
			flagDoRunMigrations,
			flagSkipMigrations,
		)
	}

	if cmd.Flags().Changed(flagDoRunMigrations) {
		run, err := cmd.Flags().GetBool(flagDoRunMigrations)
		return run, errors.Wrapf(err, "could not read flag %s", flagDoRunMigrations)
	}

	if cmd.Flags().Changed(flagSkipMigrations) {
		skip, err := cmd.Flags().GetBool(flagSkipMigrations)
		run := !skip

		return run, errors.Wrapf(err, "could not read flag %s", flagSkipMigrations)
	}

	return happyConf.AutoRunMigrations(), nil
}
