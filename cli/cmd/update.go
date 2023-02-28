package cmd

import (
	"context"

	ab "github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/config"

	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
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
		checkCreateFlags,
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric),
}

func runUpdate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	happyConfig, stackService, artifactBuilder, stackTags, awsBackend, err := initializeHappyClients(
		cmd,
		sliceName,
		tag,
		createTag,
		dryRun,
	)
	if err != nil {
		return errors.Wrap(err, "unable to initialize all the happy clients")
	}

	stackName := args[0]
	err = validateHappyEnvironment(
		ctx,
		happyConfig,
		awsBackend,
		stackService,
		stackName,
		force,
		artifactBuilder,
	)
	if err != nil {
		return errors.Wrap(err, "failed one of the happy client validations")
	}

	// 1.) if the stack doesn't exist and force flag is used, call the create function first
	stack, err := stackService.GetStack(ctx, stackName)
	if err != nil {
		if force {
			stack, err = stackService.Add(ctx, stackName, dryRun)
			if err != nil {
				return errors.Wrap(err, "unable to create the stack")
			}
		} else {
			return errors.Wrap(err, "unable to get stack")
		}
	}

	// 2.) update the existing stacks
	return updateStack(
		ctx,
		stack,
		cmd,
		stackName,
		stackTags,
		force,
		artifactBuilder,
		stackService,
		happyConfig,
		awsBackend,
	)
}

func updateStack(
	ctx context.Context,
	stack *stackservice.Stack,
	cmd *cobra.Command,
	stackName string,
	tags map[string]string,
	forceFlag bool,
	artifactBuilder ab.ArtifactBuilderIface,
	stackService *stackservice.StackService,
	happyConfig *config.HappyConfig,
	awsBackend *backend.Backend,
) error {
	// 2.) update the workspace's meta variables
	// TODO: is this used? the only thing I think some old happy environments use is the priority? I guess stack tags too
	stackMeta, err := updateStackMeta(ctx, stackName, tags, happyConfig, stackService)
	if err != nil {
		return errors.Wrap(err, "unable to update the stack's meta information")
	}

	// 3.) apply the terraform for the stack
	stack = stack.WithMeta(stackMeta)
	err = stack.Apply(ctx, makeWaitOptions(stackName, awsBackend), dryRun)
	if err != nil {
		return errors.Wrap(err, "failed to apply the stack")
	}
	if dryRun {
		logrus.Debugf("cleaning up stack '%s'", stackName)
		err = stackService.Remove(ctx, stackName, false)
		if err != nil {
			return errors.Wrap(err, "unable to remove stack")
		}
	}

	// 4.) run migrations tasks
	shouldRunMigration, err := happyCmd.ShouldRunMigrations(cmd, happyConfig)
	if err != nil {
		return err
	}
	if shouldRunMigration {
		err = runMigrate(cmd, stackName)
		if err != nil {
			return errors.Wrap(err, "failed to run migrations")
		}
	}

	// 5.) print to stdout
	stack.PrintOutputs(ctx)
	return nil
}

func updateStackMeta(
	ctx context.Context,
	stackName string,
	tags map[string]string,
	happyConfig *config.HappyConfig,
	stackService *stackservice.StackService,
) (*stackservice.StackMeta, error) {
	stackMeta := stackService.NewStackMeta(stackName)
	stackMeta.Load(map[string]string{
		"happy/meta/configsecret": happyConfig.GetSecretId(),
	})
	targetBaseTag := tag
	if sliceDefaultTag != "" {
		targetBaseTag = sliceDefaultTag
	}
	err := stackMeta.Update(ctx, targetBaseTag, tags, "", stackService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update the stack meta")
	}
	return stackMeta, nil
}

/* TODO: // consolidate some stack tags
stackTags := map[string]string{}
if sliceName != "" {
	serviceImages, err := builderConfig.GetBuildServicesImage(ctx)
	if err != nil {
		return err
	}

	for service := range serviceImages {
		stackTags[service] = tag
	}
}*/
