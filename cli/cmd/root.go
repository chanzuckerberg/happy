package cmd

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/chanzuckerberg/happy/cli/cmd/hosts"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagVerbose            = "verbose"
	flagNoColor            = "no-color"
	flagDetached           = "detached"
	flagLocalstack         = "localstack"
	flagLocalstackEndpoint = "localstackendpoint"
)

var OutputFormat string = "text"
var Interactive bool = true

func init() {
	rootCmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Use this to enable verbose mode")
	rootCmd.PersistentFlags().Bool(flagNoColor, false, "Use this to disable ANSI colors")
	rootCmd.PersistentFlags().Bool(flagDetached, false, "Use this to run in detached (non-interactive) mode")
	rootCmd.PersistentFlags().Bool(flagLocalstack, false, "Use localstack mode")
	rootCmd.PersistentFlags().Bool(flagLocalstackEndpoint, false, "Localstack endpoint (defaults to (http://localhost:4566)")

	// Add nested sub-commands here
	rootCmd.AddCommand(hosts.NewHostsCommand())
}

var rootCmd = &cobra.Command{
	Use:           "happy",
	Short:         "",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		workspace_repo.StartTFCWorkerPool(cmd.Context())
		verbose, err := cmd.Flags().GetBool(flagVerbose)
		if err != nil {
			return errors.Wrap(err, "missing verbose flag")
		}
		if verbose {
			log.Info("Running in verbose mode")
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
			detached = false
		}
		Interactive = !detached

		o, _ := os.Stdout.Stat()
		if (o.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
			// Output is not a terminal, so we should not be interactive
			Interactive = false
		}

		dctx := diagnostics.BuildDiagnosticContext(cmd.Context(), Interactive)
		cmd.SetContext(dctx)

		localstackMode, err := cmd.Flags().GetBool(flagLocalstack)
		if err != nil {
			localstackMode = false
		}
		util.SetLocalstackMode(localstackMode)
		if localstackMode {
			if localstackEndpoint, err := cmd.Flags().GetString(flagLocalstackEndpoint); err == nil {
				_, err = url.ParseRequestURI(flagLocalstackEndpoint)
				if err != nil {
					return errors.Wrap(err, "localstack endpoint is not a valid url")
				}
				util.SetLocalstackEndpoint(localstackEndpoint)
			}
		}

		err = CheckLockedHappyVersion(cmd)
		if err != nil {
			return err
		}

		return errors.Wrap(err, "local environment is misconfigured")
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		WarnIfHappyOutdated(cmd)
	},
}

// Execute executes the command
func Execute() error {
	// collect the time the command was started
	ctx := context.WithValue(context.Background(), util.CmdStartContextKey, time.Now())
	dctx := diagnostics.BuildDiagnosticContext(ctx, true)
	defer diagnostics.PrintRuntimes(dctx)
	err := rootCmd.ExecuteContext(dctx)
	if err != nil {
		return err
	}
	warnings, err := dctx.GetWarnings()
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
