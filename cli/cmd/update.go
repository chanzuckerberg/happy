package cmd

import (
	"context"
	"fmt"

	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sliceDefaultTag string

func init() {
	rootCmd.AddCommand(updateCmd)
	config.ConfigureCmdWithBootstrapConfig(updateCmd)
	happyCmd.SupportUpdateSlices(updateCmd, &sliceName, &sliceDefaultTag)
	happyCmd.SetMigrationFlags(updateCmd)

	updateCmd.Flags().StringVar(&tag, "tag", "", "Tag name for docker image. Leave empty to generate one automatically.")
	updateCmd.Flags().BoolVar(&createTag, "create-tag", true, "Will build, tag, and push images when set. Otherwise, assumes images already exist.")
	updateCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
	updateCmd.Flags().BoolVar(&force, "force", false, "Force stack creation if it doesn't exist")
	updateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Prepare all resources, but do not apply any changes")
}

var updateCmd = &cobra.Command{
	Use:          "update STACK_NAME",
	Short:        "Update stack",
	Long:         "Update stack matching STACK_NAME",
	SilenceUsage: true,
	RunE:         runUpdate,
	PreRunE: happyCmd.Validate(
		happyCmd.IsTagUsedWithSkipTag,
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric),
}

func runUpdate(cmd *cobra.Command, args []string) error {
	stackName := args[0]
	happyClient, err := makeHappyClient(cmd, sliceName, stackName, []string{tag}, createTag, dryRun, ModeUpdate)
	if err != nil {
		return errors.Wrap(err, "unable to initialize the happy client")
	}

	ctx := cmd.Context()
	err = validate(
		validateConfigurationIntegirty(ctx, happyClient),
		validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
		validateTFEBackLog(ctx, dryRun, happyClient.AWSBackend),
		validateStackNameAvailable(ctx, happyClient.StackService, stackName, force),
		validateStackExistsUpdate(ctx, stackName, dryRun, happyClient),
		validateECRExists(ctx, stackName, dryRun, terraformECRTargetPathTemplate, happyClient),
		validateImageExists(ctx, createTag, skipCheckTag, happyClient.ArtifactBuilder),
	)
	if err != nil {
		return errors.Wrap(err, "failed one of the happy client validations")
	}

	// update the existing stacks
	stack, err := happyClient.StackService.GetStack(ctx, stackName)
	if err != nil {
		return errors.Wrapf(err, "stack %s doesn't exist; this should never happen", stackName)
	}
	return updateStack(ctx, cmd, stack, force, happyClient)
}

func validateStackExistsUpdate(ctx context.Context, stackName string, dryRun bool, happyClient *HappyClient) validation {
	return func() error {
		// 1.) if the stack does not exist and force flag is used, call the create function first
		_, err := happyClient.StackService.GetStack(ctx, stackName)
		if err != nil {
			if force {
				_, err = happyClient.StackService.Add(ctx, stackName, dryRun)
				if err != nil {
					return errors.Wrap(err, "unable to create the stack")
				}
			} else {
				return errors.Wrap(err, "unable to get stack")
			}
		}

		return nil
	}
}

func updateStack(ctx context.Context, cmd *cobra.Command, stack *stackservice.Stack, forceFlag bool, happyClient *HappyClient) error {
	// 1.) update the workspace's meta variables
	stackMeta, err := updateStackMeta(ctx, stack.Name, happyClient)
	if err != nil {
		return errors.Wrap(err, "unable to update the stack's meta information")
	}

	// 2.) apply the terraform for the stack
	stack = stack.WithMeta(stackMeta)
	err = stack.Apply(ctx, makeWaitOptions(stack.Name, happyClient.HappyConfig, happyClient.AWSBackend), dryRun, workspace_repo.Message(fmt.Sprintf("Happy %s Update Stack [%s]", util.GetVersion().Version, stack.Name)))
	if err != nil {
		return errors.Wrap(err, "failed to apply the stack")
	}

	if dryRun {
		if happyClient.Mode == ModeCreate {
			logrus.Debugf("cleaning up stack '%s'", stack.Name)
			err = happyClient.StackService.Remove(ctx, stack.Name, false)
			if err != nil {
				return errors.Wrap(err, "unable to remove stack")
			}
		}
		return nil
	}

	// 3.) run migrations tasks
	shouldRunMigration, err := happyCmd.ShouldRunMigrations(cmd, happyClient.HappyConfig)
	if err != nil {
		return err
	}
	if shouldRunMigration {
		err = runMigrate(cmd, stack.Name)
		if err != nil {
			return errors.Wrap(err, "failed to run migrations")
		}
	}

	// 4.) print to stdout
	stack.PrintOutputs(ctx)

	return nil
}

func updateStackMeta(ctx context.Context, stackName string, happyClient *HappyClient) (*stackservice.StackMeta, error) {
	if sliceDefaultTag != "" {
		happyClient.ArtifactBuilder.WithTags([]string{sliceDefaultTag})
	}
	// for updating and creating (unless in dry-run mode), there should only be one tag (either provided or generated)
	tag := ""
	if len(happyClient.ArtifactBuilder.GetTags()) == 1 {
		tag = happyClient.ArtifactBuilder.GetTags()[0]
	} else if !happyClient.DryRun {
		return nil, errors.New("there should only be one tag when updating or creating a stack")
	}
	username, err := happyClient.AWSBackend.GetUserName(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get calling user's name")
	}
	stackMeta := happyClient.StackService.NewStackMeta(stackName)
	return stackMeta.UpdateAll(
		tag,
		happyClient.StackTags,
		"",
		username,
		happyClient.HappyConfig.GetProjectRoot(),
		happyClient.HappyConfig,
		stackName,
		happyClient.HappyConfig.GetEnv(),
	), nil
}
