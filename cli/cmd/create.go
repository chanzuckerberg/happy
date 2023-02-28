package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	ab "github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/cli/pkg/diagnostics"
	waitoptions "github.com/chanzuckerberg/happy/cli/pkg/options"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/cli/pkg/workspace_repo"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		checkCreateFlags,
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric),
	RunE: runCreate,
}

func checkCreateFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("skip-check-tag") && !cmd.Flags().Changed("tag") {
		return errors.New("--skip-check-tag can only be used when --tag is specified")
	}

	if !createTag && !cmd.Flags().Changed("tag") {
		return errors.New("Must specify a tag when create-tag=false")
	}

	return nil
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

type validation func() error

func validateImageExists(ctx context.Context, createTag, skipCheckTag bool, ab artifact_builder.ArtifactBuilderIface) validation {
	return func() error {
		if createTag {
			// if we build and push and it succeeds, we know that the image exists
			return ab.BuildAndPush(ctx)
		}

		if skipCheckTag {
			return nil
		}

		if len(ab.GetTags()) == 0 {
			return errors.Errorf("no tags have been assigned")
		}

		exists, err := ab.CheckImageExists(ctx, ab.GetTags()[0])
		if err != nil {
			return errors.Wrapf(err, "error checking if tag %s existed", ab.GetTags()[0])
		}
		if !exists {
			return errors.Errorf("image tag does not exist: '%s'", ab.GetTags()[0])
		}

		return nil
	}
}
func validateTFEBackLog(ctx context.Context, isDryRun bool, awsBackend *backend.Backend) validation {
	return func() error {
		return verifyTFEBacklog(ctx, createWorkspaceRepo(isDryRun, awsBackend))
	}
}

func validateGitTree(projectRoot string) validation {
	return func() error {
		return util.ValidateGitTree(projectRoot)
	}
}

func validateStackNameAvailable(ctx context.Context, stackService *stackservice.StackService, stackName string, force bool) validation {
	return func() error {
		if force {
			return nil
		}

		_, err := stackService.GetStack(ctx, stackName)
		if err != nil {
			return nil
		}

		return errors.Wrap(err, "the stack name is already taken")
	}
}

func validate(validations ...validation) error {
	for _, validation := range validations {
		err := validation()
		if err != nil {
			return errors.Wrap(err, "unable to validate the environment")
		}
	}
	return nil
}

func makeWaitOptions(stackName string, backend *backend.Backend) waitoptions.WaitOptions {
	taskOrchestrator := orchestrator.NewOrchestrator().WithBackend(backend)
	return waitoptions.WaitOptions{
		StackName:    stackName,
		Orchestrator: taskOrchestrator,
		Services:     backend.Conf().GetServices(),
	}
}

func verifyTFEBacklog(ctx context.Context, workspaceRepo workspace_repo.WorkspaceRepoIface) error {
	if !diagnostics.IsInteractiveContext(ctx) {
		// When you're not interactive, no point in measuring the backlog size
		return nil
	}
	backlogSize, _, err := workspaceRepo.EstimateBacklogSize(ctx)
	if err != nil {
		return errors.Wrap(err, "error estimating TFE backlog")
	}
	if backlogSize < 2 {
		logrus.Debug("There is no TFE backlog, proceeding.")
	} else if backlogSize < 20 {
		logrus.Debugf("TFE backlog is only %d runs long, proceeding.", backlogSize)
	} else {
		proceed := false
		prompt := &survey.Confirm{Message: fmt.Sprintf("TFE backlog is %d runs long, it might take a while to clear out. Do you want to wait? ", backlogSize)}
		err = survey.AskOne(prompt, &proceed)
		if err != nil {
			return errors.Wrapf(err, "failed to ask for confirmation")
		}

		if !proceed {
			return err
		}
	}
	return nil
}

func validateHappyEnvironment(
	ctx context.Context,
	happyConfig *config.HappyConfig,
	awsBackend *backend.Backend,
	stackService *stackservice.StackService,
	stackName string,
	force bool,
	artifactBuilder ab.ArtifactBuilderIface,
) error {
	return validate(
		validateGitTree(happyConfig.GetProjectRoot()),
		validateTFEBackLog(ctx, dryRun, awsBackend),
		validateStackNameAvailable(ctx, stackService, stackName, force),
		validateImageExists(ctx, createTag, skipCheckTag, artifactBuilder),
	)
}
