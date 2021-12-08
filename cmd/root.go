package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

const (
	flagVerbose = "verbose"
)

func init() {
	rootCmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Use this to enable verbose mode")
}

var rootCmd = &cobra.Command{
	Use:   "happy",
	Short: "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		verbose, err := cmd.Flags().GetBool(flagVerbose)
		if err != nil {
			return errors.Wrap(err, "Missing verbose flag")
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
		}
		return nil
	},
}

// Execute executes the command
func Execute() error {
	return rootCmd.Execute()
}
