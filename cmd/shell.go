package cmd

import (
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shellCmd)
	config.ConfigureCmdWithBootstrapConfig(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:          "shell STACK_NAME SERVICE",
	Short:        "",
	Long:         "",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(2)),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

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

		b, err := backend.NewAWSBackend(ctx, happyConfig)
		if err != nil {
			return err
		}

		return orchestrator.NewOrchestrator().WithBackend(b).Shell(stackName, service)
	},
}
