package cmd

import (
	"errors"
	"os"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell STACK_NAME SERVICE",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		env := "rdev"

		if len(args) != 2 {
			return errors.New("Incorrect number of arguments")
		}

		stackName := args[0]
		service := args[1]

		happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
		if !ok {
			return errors.New("Please set env var HAPPY_CONFIG_PATH")
		}

		_, ok = os.LookupEnv("HAPPY_PROJECT_ROOT")
		if !ok {
			return errors.New("Please set env var HAPPY_PROJECT_ROOT")
		}

		happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
		if err != nil {
			return err
		}

		taskRunner := backend.GetAwsEcs(happyConfig)
		taskOrchestrator := orchestrator.NewOrchestrator(happyConfig, taskRunner)
		err = taskOrchestrator.Shell(stackName, service)

		return err
	},
}
