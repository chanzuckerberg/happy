package cmd

import (
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
	rootCmd.AddCommand(migrateCmd)
	config.ConfigureCmdWithBootstrapConfig(migrateCmd)

	migrateCmd.Flags().StringSliceVar(&pushImages, "images", []string{}, "List of images to migrate to registry.")
	migrateCmd.Flags().BoolVar(&reset, "reset", false, "Resetting the task")
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

	url := happyConfig.TfeUrl()
	org := happyConfig.TfeOrg()

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
