package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chanzuckerberg/happy/cmd/hosts"
	"github.com/chanzuckerberg/happy/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagVerbose  = "verbose"
	flagNoColor  = "no-color"
	flagDetached = "detached"
)

var Interactive bool = true

func init() {
	rootCmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Use this to enable verbose mode")
	rootCmd.PersistentFlags().Bool(flagNoColor, false, "Use this to disable ANSI colors")
	rootCmd.PersistentFlags().Bool(flagDetached, false, "Use this to run in detached (non-interactive) mode")

	// Add nested sub-commands here
	rootCmd.AddCommand(hosts.NewHostsCommand())
}

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

		detached, err := cmd.Flags().GetBool(flagDetached)
		if err != nil {
			return errors.Wrap(err, "missing detached flag")
		}
		Interactive = !detached

		err = util.ValidateEnvironment(context.Background())
		return errors.Wrap(err, "local environment is misconfigured")
	},
}

// Execute executes the command
func Execute() error {
	ctx := diagnostics.BuildDiagnosticContext(context.Background())
	defer diagnostics.PrintRuntimes(ctx)
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		return err
	}
	warnings, err := ctx.GetWarnings()
	if err != nil {
		return errors.Wrap(err, "failed to get warnings")
	}
	if len(warnings) > 0 {
		log.Warn("Warnings:")
		for _, warning := range warnings {
			log.Warn(warning)
		}
	}
	return nil
}

func PrintError(err error) {
	os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
}

func printOutput(output string) {
	os.Stdout.WriteString(fmt.Sprintf("%s\n", output))
}
