package cmd

import (
	"os"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	stack_service "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var reset bool

func init() {
	migrateCmd.Flags().StringSliceVar(&pushImages, "images", []string{}, "List of images to migrate to registry.")
	migrateCmd.Flags().BoolVar(&reset, "reset", false, "Resetting the task")
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate STACK_NAME",
	Short: "migrate stack",
	Long:  "Run migration tasks for stack with given name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]

		return runMigrate(stackName)
	},
}

func runMigrate(stackName string) error {

	env := "rdev"

	happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
	if !ok {
		return errors.New("please set env var HAPPY_CONFIG_PATH")
	}

	_, ok = os.LookupEnv("HAPPY_PROJECT_ROOT")
	if !ok {
		return errors.New("please set env var HAPPY_PROJECT_ROOT")
	}

	happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
	if err != nil {
		return err
	}
	taskRunner := backend.GetAwsEcs(happyConfig)
	taskOrchestrator := orchestrator.NewOrchestrator(happyConfig, taskRunner)

	url, err := happyConfig.TfeUrl()
	if err != nil {
		return err
	}
	org, err := happyConfig.TfeOrg()
	if err != nil {
		return err
	}
	workspaceRepo, err := workspace_repo.NewWorkspaceRepo(url, org)
	if err != nil {
		return err
	}
	paramStoreBackend := backend.GetAwsBackend(happyConfig)
	stackService := stack_service.NewStackService(happyConfig, paramStoreBackend, workspaceRepo)

	stacks, err := stackService.GetStacks()
	if err != nil {
		return err
	}
	stack, ok := stacks[stackName]
	if !ok {
		return errors.Errorf("stack %s not found", stackName)
	}

	showLogs := true
	if reset {
		err = taskOrchestrator.RunTasks(stack, string(backend.DeletionTask), showLogs)
		if err != nil {
			return err
		}
	}

	return taskOrchestrator.RunTasks(stack, string(backend.MigrationTask), showLogs)
}
