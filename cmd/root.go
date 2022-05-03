package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/cmd/hosts"
	"github.com/chanzuckerberg/happy/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagVerbose = "verbose"
	flagNoColor = "no-color"
)

func init() {
	rootCmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Use this to enable verbose mode")
	rootCmd.PersistentFlags().Bool(flagNoColor, false, "Use this to disable ANSI colors")

	// Add nested sub-commands here
	rootCmd.AddCommand(hosts.NewHostsCommand())
}

var ctx = diagnostics.BuildDiagnosticContext(context.Background())

var rootCmd = &cobra.Command{
	Use:           "happy",
	Short:         "",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		verbose, err := cmd.Flags().GetBool(flagVerbose)
		if err != nil {
			return errors.Wrap(err, "missing verbose flag")
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
		}

		noColor, err := cmd.Flags().GetBool(flagNoColor)
		if err != nil {
			return errors.Wrap(err, "missing no-color flag")
		}
		if noColor {
			color.NoColor = noColor
		}

		err = util.ValidateEnvironment(context.Background())
		return errors.Wrap(err, "local environment is misconfigured")
	},
}

// Execute executes the command
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	warnings := ctx.GetWarnings()
	if len(warnings) > 0 {
		log.Warn("Warnings:")
		for _, warning := range warnings {
			log.Warn(warning)
		}
	}
	return nil
}
