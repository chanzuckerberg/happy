package cmd

import (
	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
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
		happyCmd.IsTagUsedWithSkipTag(createTag),
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric),
	RunE: runCreate,
}

func runCreate(
	cmd *cobra.Command,
	args []string,
) error {
	ctx := cmd.Context()
	happyConfig, stackService, artifactBuilder, stackTags, awsBackend, err := initializeHappyClients(
		cmd,
		sliceName,
		tag,
		createTag,
		dryRun,
	)
	if err != nil {
		return err
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

	// 1.) if the stack does not exist and force flag is used, call the create function first
	stack, err := stackService.GetStack(ctx, stackName)
	if err != nil {
		stack, err = stackService.Add(ctx, stackName, dryRun)
		if err != nil {
			return errors.Wrap(err, "unable to create the stack")
		}
	} else {
		if !force {
			return errors.Wrapf(err, "stack %s already exists", stackName)
		}
	}

	// 2.) otherwise, update the existing stacks
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
