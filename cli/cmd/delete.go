package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/options"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running delete")
		for _, stackName := range args {
			err := runDelete(cmd, stackName)
			if err != nil {
				return err
			}
		}
		return nil
	},
	PreRunE: happyCmd.Validate(
		cobra.MinimumNArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric,
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()
			return util.ValidateEnvironment(cmd.Context(),
				checklist.TerraformInstalled,
			)
		},
	),
}

func runDelete(cmd *cobra.Command, stackName string) error {
	happyClient, err := makeHappyClient(cmd, sliceName, stackName, []string{tag}, createTag)
	if err != nil {
		return errors.Wrap(err, "unable to initialize the happy client")
	}
	ctx := context.WithValue(cmd.Context(), options.DryRunKey, dryRun)
	message := workspace_repo.Message(fmt.Sprintf("Happy %s Delete Stack [%s]", util.GetVersion().Version, stackName))
	err = validate(
		validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
		validateStackNameAvailable(ctx, happyClient.StackService, stackName, force),
		validateStackExists(ctx, stackName, happyClient, message),
		validateTFEBackLog(ctx, happyClient.AWSBackend),
	)
	if err != nil {
		log.Warnf("failed one of the happy client validations %s", err.Error())
		return err
	}

	stacks, err := happyClient.StackService.GetStacks(ctx)
	if err != nil {
		return err
	}

	stack, ok := stacks[stackName]
	if !ok {
		return errors.Wrapf(err, "stack %s not found", stackName)
	}

	// Run all necessary tasks before deletion
	taskOrchestrator := orchestrator.
		NewOrchestrator().
		WithHappyConfig(happyClient.HappyConfig).
		WithBackend(happyClient.AWSBackend)
	err = taskOrchestrator.RunTasks(ctx, stack, backend.TaskTypeDelete)
	if err != nil {
		if !force {
			if !diagnostics.IsInteractiveContext(ctx) {
				return err
			}
			proceed := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Error running tasks while trying to delete %s (%s); Continue? ", stackName, err.Error()),
			}
			err = survey.AskOne(prompt, &proceed)
			if err != nil {
				return errors.Wrapf(err, "failed to ask for confirmation")
			}
			if !proceed {
				return err
			}
		}
	}

	hasState, err := happyClient.StackService.HasState(ctx, stackName)
	if err != nil {
		return errors.Wrapf(err, "unable to determine whether the stack has state")
	}

	runopts := workspace_repo.Message(fmt.Sprintf("Happy %s Delete Stack [%s]", util.GetVersion().Version, stackName))
	if !hasState {
		log.Info("No state found for stack, workspace will be removed")
		return happyClient.StackService.Remove(ctx, stackName, runopts)
	}

	// Destroy the stack
	destroySuccess := true
	waitopts := options.WaitOptions{
		StackName:    stackName,
		Orchestrator: taskOrchestrator,
		Services:     happyClient.HappyConfig.GetServices(),
	}

	tfDirPath := happyClient.HappyConfig.TerraformDirectory()
	happyProjectRoot := happyClient.HappyConfig.GetProjectRoot()
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)
	if err = stack.Destroy(ctx, srcDir, waitopts, runopts); err != nil {
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
		return happyClient.StackService.Remove(ctx, stackName, runopts)
	}
	log.Warnf("Stack %s was not deleted fully", stackName)
	return nil
}
