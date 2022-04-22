package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	config.ConfigureCmdWithBootstrapConfig(deleteCmd)

	deleteCmd.Flags().BoolVar(&force, "force", false, "Force stack deletion")
}

var deleteCmd = &cobra.Command{
	Use:          "delete STACK_NAME",
	Short:        "delete an existing stack",
	Long:         "Delete the stack with the given name.",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1), cmd.CheckStackName),
	RunE:         runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	stackName := args[0]

	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
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
		if !force {
			proceed := false
			prompt := &survey.Confirm{Message: fmt.Sprintf("Error running tasks while trying to delete %s (%s); Continue? ", stackName, err.Error())}
			err = survey.AskOne(prompt, &proceed)
			if err != nil {
				return errors.Wrapf(err, "failed to ask for confirmation")
			}
			if !proceed {
				return err
			}
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
		hasState, err := stackService.HasState(ctx, stackName)
		if err != nil {
			return errors.Wrapf(err, "unable to determine whether the stack has state")
		}
		if hasState {
			proceed := false
			prompt := &survey.Confirm{Message: fmt.Sprintf("Error while destroying %s; resources might remain. Continue to remove workspace? ", stackName)}
			err = survey.AskOne(prompt, &proceed)
			if err != nil {
				return errors.Wrapf(err, "failed to ask for confirmation")
			}

			if !proceed {
				return err
			}
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
