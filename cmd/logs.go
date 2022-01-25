package cmd

import (
	"errors"
	"os"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	"github.com/spf13/cobra"
)

var since string

func init() {
	logsCmd.Flags().StringVar(&since, "since", "10m", "Length of time to look back in logs")
	rootCmd.AddCommand(logsCmd)
}

var logsCmd = &cobra.Command{
	Use:   "logs STACK_NAME SERVICE",
	Short: "Tail logs",
	Long:  "Tail the logs of a service (frontend, backend, upload, migrations)",
	RunE:  runLogs,
}

func runLogs(cmd *cobra.Command, args []string) error {

	env := "rdev"

	if len(args) != 2 {
		return errors.New("incorrect number of arguments")
	}

	stackName := args[0]
	service := args[1]

	happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
	if !ok {
		return errors.New("please set env var HAPPY_CONFIG_PATH")
	}

	happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
	if err != nil {
		return err
	}
	taskRunner := backend.GetAwsEcs(happyConfig)
	taskOrchestrator := orchestrator.NewOrchestrator(happyConfig, taskRunner)

	return taskOrchestrator.Logs(stackName, service, since)
}
