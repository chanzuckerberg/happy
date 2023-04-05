package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	config.ConfigureCmdWithBootstrapConfig(deleteCmd)

	deleteCmd.Flags().BoolVar(&force, "force", false, "Force stack deletion")
	deleteCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Prepare all resources, but do not apply any changes")
}

var deleteCmd = &cobra.Command{
	Use:          "delete STACK_NAME",
	Short:        "Delete an existing stack",
	Long:         "Delete the stack with the given name.",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1), cmd.IsStackNameDNSCharset),
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

	b, err := backend.NewAWSBackend(ctx, happyConfig.GetEnvironmentContext())
	if err != nil {
		return err
	}

	workspaceRepo := createWorkspaceRepo(dryRun, b)
	stackService := stackservice.NewStackService().WithHappyConfig(happyConfig).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

	err = verifyTFEBacklog(ctx, workspaceRepo)
	if err != nil {
		return err
	}

	// FIXME TODO check env to make sure it allows for stack deletion
	if dryRun {
		log.Debugf("Planning removal of stack '%s'\n", stackName)
	} else {
		log.Debugf("Deleting stack '%s'\n", stackName)
	}
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
	taskOrchestrator := orchestrator.NewOrchestrator().WithHappyConfig(happyConfig).WithBackend(b).WithDryRun(dryRun)
	err = taskOrchestrator.RunTasks(ctx, stack, backend.TaskTypeDelete)
	if err != nil {
		if !force {
			if !diagnostics.IsInteractiveContext(ctx) {
				return err
			}
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

	hasState, err := stackService.HasState(ctx, stackName)
	if err != nil {
		return errors.Wrapf(err, "unable to determine whether the stack has state")
	}

	if !hasState {
		log.Info("No state found for stack, workspace will be removed")
		return removeWorkspace(ctx, stackService, stackName, dryRun)
	}

	options := workspace_repo.Message(fmt.Sprintf("Happy %s Delete Stack [%s]", util.GetVersion().Version, stackName))

	// Destroy the stack
	destroySuccess := true
	if err = stack.PlanDestroy(ctx, dryRun, options); err != nil {
		// log error and set a flag, but do not return
		log.Errorf("Failed to destroy stack: '%s'", err)
		destroySuccess = false
	}

	doRemoveWorkspace := false
	if !destroySuccess {
		if !diagnostics.IsInteractiveContext(ctx) {
			return errors.Errorf("Error while destroying %s; resources might remain, aborting workspace removal in non-interactive mode.", stackName)
		}

		proceed := false
		prompt := &survey.Confirm{Message: fmt.Sprintf("Error while destroying %s; resources might remain. Continue to remove workspace? ", stackName)}
		err = survey.AskOne(prompt, &proceed)
		if err != nil {
			return errors.Wrapf(err, "failed to ask for confirmation")
		}

		if !proceed {
			return err
		}

		doRemoveWorkspace = true
	}

	// Remove the stack from state
	// TODO: are these the right error messages?
	if destroySuccess || doRemoveWorkspace {
		return removeWorkspace(ctx, stackService, stackName, dryRun, options)
	} else {
		log.Println("Delete NOT done")
	}

	return nil
}

func removeWorkspace(ctx context.Context, stackService *stackservice.StackService, stackName string, dryRun bool, options ...workspace_repo.TFERunOption) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "removeWorkspace")
	err := stackService.Remove(ctx, stackName, dryRun, options...)
	if err != nil {
		return err
	}
	log.Println("Delete done")
	return nil
}
