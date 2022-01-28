package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagVerbose = "verbose"

	flagHappyProjectRoot        = "happy-project-root"
	flagHappyConfigPath         = "happy-config-path"
	flagDockerComposeConfigPath = "docker-compose-config-path"
)

// We will load bootrap configuration common to all commands here
// can then be consumed by other commands as needed.
var (
	bootstrapConfig         *config.Bootstrap
	happyProjectRoot        string
	happyConfigPath         string
	dockerComposeConfigPath string
)

func init() {
	rootCmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Use this to enable verbose mode")

	rootCmd.PersistentFlags().StringVar(&happyProjectRoot, flagHappyProjectRoot, "", "Specify the root of your Happy project")
	rootCmd.PersistentFlags().StringVar(&happyConfigPath, flagHappyConfigPath, "", "Specify the path to your Happy project's config file")
	rootCmd.PersistentFlags().StringVar(&dockerComposeConfigPath, flagDockerComposeConfigPath, "", "Specify the path to your Happy project's docker compose file")
}

var rootCmd = &cobra.Command{
	Use:   "happy",
	Short: "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		verbose, err := cmd.Flags().GetBool(flagVerbose)
		if err != nil {
			return errors.Wrap(err, "missing verbose flag")
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
		}

		bootstrapConfig, err = config.ResolveBootstrapConfig()
		return err
	},
}

// Execute executes the command
func Execute() error {
	return rootCmd.Execute()
}
