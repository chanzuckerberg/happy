package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/options"
	stackservice "github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	sliceDefaultTag string
)

func init() {
	rootCmd.AddCommand(updateCmd)
	config.ConfigureCmdWithBootstrapConfig(updateCmd)
	happyCmd.SupportUpdateSlices(updateCmd, &sliceName, &sliceDefaultTag)
	happyCmd.SetMigrationFlags(updateCmd)
	happyCmd.SetImagePromotionFlags(updateCmd, &imageSrcEnv, &imageSrcStack, &imageSrcRoleArn)
	happyCmd.SetDryRunFlag(updateCmd, &dryRun)

	updateCmd.Flags().StringVar(&tag, "tag", "", "Tag name for docker image. Leave empty to generate one automatically.")
	updateCmd.Flags().BoolVar(&createTag, "create-tag", true, "Will build, tag, and push images when set. Otherwise, assumes images already exist.")
	updateCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
	updateCmd.Flags().BoolVar(&force, "force", false, "Force stack creation if it doesn't exist")
}

var updateCmd = &cobra.Command{
	Use:          "update STACK_NAME",
	Short:        "Update stack",
	Long:         "Update stack matching STACK_NAME",
	SilenceUsage: true,
	RunE:         runUpdate,
	PreRunE: happyCmd.Validate(
		happyCmd.IsImageEnvUsedWithImageStack,
		happyCmd.IsTagUsedWithSkipTag,
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric,
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()

			// Required for all commmands
			required_checks := []util.ValidationCallback{
				checklist.TerraformInstalled,
				checklist.AwsInstalled,
			}

			if !skipCheckTag || createTag {
				required_checks = append(required_checks, checklist.MinDockerComposeVersion, checklist.DockerEngineRunning, checklist.DockerInstalled)
			}

			return util.ValidateEnvironment(cmd.Context(), required_checks...)
		},
	),
}

func runUpdate(cmd *cobra.Command, args []string) error {
	stackName := args[0]
	happyClient, err := makeHappyClient(cmd, sliceName, stackName, []string{tag}, createTag)
	if err != nil {
		return errors.Wrap(err, "unable to initialize the happy client")
	}

	ctx := context.WithValue(cmd.Context(), options.DryRunKey, dryRun)
	err = validate(
		validateConfigurationIntegirty(ctx, sliceName, happyClient),
		validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
		validateStackNameAvailable(ctx, happyClient.StackService, stackName, force),
		validateTFEBackLog(ctx, happyClient.AWSBackend),
		validateStackExistsUpdate(ctx, stackName, happyClient),
		validateECRExists(ctx, stackName, terraformECRTargetPathTemplate, happyClient),
		validateImageExists(ctx, createTag, skipCheckTag, imageSrcEnv, imageSrcStack, imageSrcRoleArn, happyClient, cmd.Flags().Changed(config.FlagAWSProfile)),
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

func validateStackExistsUpdate(ctx context.Context, stackName string, happyClient *HappyClient) validation {
	return func() error {
		// 1.) if the stack does not exist and force flag is used, call the create function first
		_, err := happyClient.StackService.GetStack(ctx, stackName)
		if err != nil {
			if force {
				_, err = happyClient.StackService.Add(ctx, stackName)
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

func validateStackExists(ctx context.Context, stackName string, happyClient *HappyClient, options ...workspace_repo.TFERunOption) validation {
	log.Debug("Scheduling validateStackExists()")
	return func() error {
		log.Debug("Running validateStackExists()")
		_, err := happyClient.StackService.GetStack(ctx, stackName)
		if err != nil {
			return errors.Wrapf(err, "stack %s doesn't exist", stackName)
		}
		return nil
	}
}

func updateStack(ctx context.Context, cmd *cobra.Command, stack *stackservice.Stack, forceFlag bool, happyClient *HappyClient) error {

	stackInfo, err := stack.GetStackInfo(ctx)

	// 1.) update the workspace's meta variables
	stackMeta, err := updateStackMeta(ctx, stack.Name, happyClient)
	if err != nil {
		return errors.Wrap(err, "unable to update the stack's meta information")
	}

	// 2.) apply the terraform for the stack
	stack = stack.WithMeta(stackMeta)
	tfDirPath := happyClient.HappyConfig.TerraformDirectory()
	happyProjectRoot := happyClient.HappyConfig.GetProjectRoot()
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)
	err = stack.Apply(ctx, srcDir, makeWaitOptions(stack.Name, happyClient.HappyConfig, happyClient.AWSBackend), workspace_repo.Message(fmt.Sprintf("Happy %s Update Stack [%s]", util.GetVersion().Version, stack.Name)))
	if err != nil {
		return errors.Wrap(err, "failed to apply the stack")
	}

	// 3.) run migrations tasks
	shouldRunMigration, err := happyCmd.ShouldRunMigrations(ctx, cmd, happyClient.HappyConfig)
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

	// Remove images with the previous tag from all ECRs, unless the previous tag is the same as the current tag
	found := false
	for _, tag := range happyClient.ArtifactBuilder.GetTags() {
		if tag == stackInfo.StackMetadata.Tag {
			found = true
			break
		}
	}

	if !found {
		err = happyClient.ArtifactBuilder.DeleteImages(ctx, stackInfo.StackMetadata.Tag)
		if err != nil {
			return errors.Wrap(err, "failed to delete images")
		}
	}

	return nil
}

func updateStackMeta(ctx context.Context, stackName string, happyClient *HappyClient) (*stackservice.StackMeta, error) {
	if sliceDefaultTag != "" {
		happyClient.ArtifactBuilder.WithTags([]string{sliceDefaultTag})
	}
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	// for updating and creating (unless in dry-run mode), there should only be one tag (either provided or generated)
	tag := ""
	if len(happyClient.ArtifactBuilder.GetTags()) == 1 {
		tag = happyClient.ArtifactBuilder.GetTags()[0]
	} else if !dryRun {
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
