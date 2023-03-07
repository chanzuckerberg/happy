package cmd

import (
	"context"
	"fmt"

	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/cli/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	force        bool
	skipCheckTag bool
	createTag    bool
	tag          string
	dryRun       bool
)

func init() {
	rootCmd.AddCommand(createCmd)
	config.ConfigureCmdWithBootstrapConfig(createCmd)
	happyCmd.SupportUpdateSlices(createCmd, &sliceName, &sliceDefaultTag) // Should this function be renamed to something more generalized?
	happyCmd.SetMigrationFlags(createCmd)

	createCmd.Flags().StringVar(&tag, "tag", "", "Specify the tag for the docker images. If not specified we will generate a default tag.")
	createCmd.Flags().BoolVar(&createTag, "create-tag", true, "Will build, tag, and push images when set. Otherwise, assumes images already exist.")
	createCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
	createCmd.Flags().BoolVar(&force, "force", false, "Ignore the already-exists errors")
	createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Plan all infrastructure changes, but do not apply them")
}

var createCmd = &cobra.Command{
	Use:          "create STACK_NAME",
	Short:        "Create new stack",
	Long:         "Create a new stack with a given tag.",
	SilenceUsage: true,
	PreRunE: happyCmd.Validate(
		happyCmd.IsTagUsedWithSkipTag,
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric),
	RunE: runCreate,
}

// keep in sync with happy-stack-eks terraform module
const terraformECRTargetPathTemplate = `module.stack.module.services["%s"].module.ecr`

func runCreate(
	cmd *cobra.Command,
	args []string,
) error {
	stackName := args[0]
	happyClient, err := makeHappyClient(cmd, sliceName, stackName, tag, createTag, dryRun)
	if err != nil {
		return errors.Wrap(err, "unable to initialize the happy client")
	}

	ctx := cmd.Context()
	err = validate(
		validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
		validateTFEBackLog(ctx, dryRun, happyClient.AWSBackend),
		validateStackNameAvailable(ctx, happyClient.StackService, stackName, force),
		validateStackExistsCreate(ctx, stackName, dryRun, happyClient),
		validateECRExists(ctx, stackName, dryRun, terraformECRTargetPathTemplate, happyClient),
		validateImageExists(ctx, createTag, skipCheckTag, happyClient.ArtifactBuilder),
	)
	if err != nil {
		return errors.Wrap(err, "failed one of the happy client validations")
	}

	// update the newly created stack
	stack, err := happyClient.StackService.GetStack(ctx, stackName)
	if err != nil {
		return errors.Wrapf(err, "stack %s doesn't exist; this should never happen", stackName)
	}
	return updateStack(ctx, cmd, stack, force, happyClient)
}

func validateECRExists(ctx context.Context, stackName string, dryRun bool, ecrTargetPathFormat string, happyClient *HappyClient) validation {
	return func() error {
		if !happyClient.HappyConfig.GetFeatures().EnableECRAutoCreation {
			return nil
		}

		targetAddrs := []string{}
		for _, service := range happyClient.HappyConfig.GetServices() {
			targetAddrs = append(targetAddrs, fmt.Sprintf(ecrTargetPathFormat, service))
		}
		stack, err := happyClient.StackService.GetStack(ctx, stackName)
		if err != nil {
			return errors.Wrapf(err, "stack %s doesn't exist; this should never happen", stackName)
		}
		stackMeta, err := updateStackMeta(ctx, stack.Name, happyClient)
		if err != nil {
			return errors.Wrap(err, "unable to update the stack's meta information")
		}

		// this has a strong coupling with the TF version that we are using,
		// so if the user isn't on it yet, this will fail
		// TODO: maybe CDK
		stack = stack.WithMeta(stackMeta)
		return stack.Apply(ctx, makeWaitOptions(stackName, happyClient.AWSBackend), dryRun, workspace_repo.TargetAddrs(targetAddrs))
	}
}

func validateStackExistsCreate(ctx context.Context, stackName string, dryRun bool, happyClient *HappyClient) validation {
	return func() error {
		// 1.) if the stack does not exist and force flag is used, call the create function first
		_, err := happyClient.StackService.GetStack(ctx, stackName)
		if err != nil {
			_, err = happyClient.StackService.Add(ctx, stackName, dryRun)
			if err != nil {
				return errors.Wrap(err, "unable to create the stack")
			}
		} else {
			if !force {
				return errors.Wrapf(err, "stack %s already exists", stackName)
			}
		}

		return nil
	}
}
