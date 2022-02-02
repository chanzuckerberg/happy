package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shellCmd)
	config.ConfigureCmdWithBootstrapConfig(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell STACK_NAME SERVICE",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		service := args[1]

		bootstrapConfig, err := config.NewBootstrapConfig()
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(bootstrapConfig)
		if err != nil {
			return err
		}

		taskRunner := backend.GetAwsEcs(happyConfig)
		taskOrchestrator := orchestrator.NewOrchestrator(happyConfig, taskRunner)
		return taskOrchestrator.Shell(stackName, service)
	},
}
