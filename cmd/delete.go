package cmd

import (
	"fmt"

	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	config.ConfigureCmdWithBootstrapConfig(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:          "delete STACK_NAME",
	Short:        "delete an existing stack",
	Long:         "Delete the stack with the given name.",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1)),
	RunE:         runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	stackName := args[0]

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

	url := b.Conf().GetTfeUrl()
	org := b.Conf().GetTfeOrg()

	workspaceRepo := workspace_repo.NewWorkspaceRepo(url, org)
	if err != nil {
		return err
	}

	stackService := stackservice.NewStackService().WithBackend(b).WithWorkspaceRepo(workspaceRepo)

	// FIXME TODO check env to make sure it allows for stack deletion
	log.Infof("Deleting %s\n", stackName)
	stacks, err := stackService.GetStacks(ctx)
	if err != nil {
		return err
	}

	// if stack not found, we're done
	stack, ok := stacks[stackName]
	if !ok {
		log.Infof("stack %s not found, no further action", stackName)
		return nil
	}

	// Run all necessary tasks before deletion
	taskOrchestrator := orchestrator.NewOrchestrator().WithBackend(b)
	err = taskOrchestrator.RunTasks(ctx, stack, string(backend.TaskTypeDelete))
	if err != nil {
		log.Errorf("Error running tasks while trying to delete %s (%s); Continue (y/n)? ", stackName, err.Error())
		var ans string
		fmt.Scanln(&ans)
		YES := map[string]bool{"Y": true, "y": true, "yes": true, "YES": true}
		if _, ok := YES[ans]; !ok {
			return err
		}
	}

	// Destroy the stack
	destroySuccess := true
	if err = stack.Destroy(ctx); err != nil {
		// log error and set a flag, but do not return
		log.Infof("Failed to destroy stack %s", err)
		destroySuccess = false
	}

	doRemoveWorkspace := false
	if !destroySuccess {
		log.Infof("Error while destroying %s; resources might remain. Continue to remove workspace (y/n)? ", stackName)
		var ans string
		fmt.Scanln(&ans)
		YES := map[string]bool{"Y": true, "y": true, "yes": true, "YES": true}
		if _, ok := YES[ans]; ok {
			doRemoveWorkspace = true
		}
	}

	// Remove the stack from state
	// TODO: are these the right error messages?
	if destroySuccess || doRemoveWorkspace {
		err = stackService.Remove(ctx, stackName)
		if err != nil {
			return err
		}
		log.Println("Delete done")
	} else {
		log.Println("Delete NOT done")
	}

	return nil
}
