package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	stack_service "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	config.ConfigureCmdWithBootstrapConfig(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete STACK_NAME",
	Short: "delete an existing stack",
	Long:  "Delete the stack with the given name.",
	RunE:  runDelete,
	Args:  cobra.ExactArgs(1),
}

func runDelete(cmd *cobra.Command, args []string) error {
	stackName := args[0]

	bootstrapConfig, err := config.NewBootstrapConfig()
	if err != nil {
		return err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return err
	}

	taskRunner := backend.GetAwsEcs(happyConfig)

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

	// TODO check env to make sure it allows for stack deletion
	log.Printf("Deleting %s\n", stackName)
	stacks, err := stackService.GetStacks()
	if err != nil {
		return err
	}

	stack, ok := stacks[stackName]
	if !ok {
		return errors.Errorf("stack %s not found", stackName)
	}

	// Run all necessary tasks before deletion
	showLogs := true
	taskOrchestrator := orchestrator.NewOrchestrator(happyConfig, taskRunner)
	err = taskOrchestrator.RunTasks(stack, string(backend.DeletionTask), showLogs)
	if err != nil {
		var ans string
		fmt.Scanln(&ans)
		log.Printf("Error running tasks while trying to delete %s; Continue (y/n)? ", stackName)
		YES := map[string]bool{"Y": true, "y": true, "yes": true, "YES": true}
		if _, ok := YES[ans]; !ok {
			return err
		}
	}

	// Destroy the stack
	destroySuccess := true
	if err = stack.Destroy(); err != nil {
		// log error and set a flag, but not do not return
		log.Printf("Failed to destroy stack %s", err)
		destroySuccess = false
	}

	doRemoveWorkspace := false
	if !destroySuccess {
		var ans string
		fmt.Scanln(&ans)
		log.Printf("Error while destroying %s; resources might remain. Continue to remove workspace (y/n)? ", stackName)
		YES := map[string]bool{"Y": true, "y": true, "yes": true, "YES": true}
		if _, ok := YES[ans]; ok {
			doRemoveWorkspace = true
		}
	}

	// Remove the stack from state
	if destroySuccess || doRemoveWorkspace {
		err = stackService.Remove(stackName)
		if err != nil {
			return err
		}
		log.Println("Delete done")
	} else {
		log.Println("Delete NOT done")
	}

	return nil
}
