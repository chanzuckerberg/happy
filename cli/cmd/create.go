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
	happyConfig, stackService, artifactBuilder, awsBackend, err := initializeHappyClients(cmd, sliceName, tag, createTag, dryRun)
	if err != nil {
		return err
	}

	stackName := args[0]
	err = validate(
		validateGitTree(happyConfig.GetProjectRoot()),
		validateTFEBackLog(ctx, dryRun, awsBackend),
		validateExistingStack(ctx, stackService, stackName, force),
		validateImageExists(ctx, skipCheckTag, artifactBuilder),
	)
	if err != nil {
		return errors.Wrap(err, "failed one of client validations")
	}

	return createStack(
		ctx,
		cmd,
		stackName,
		map[string]string{},
		force,
		artifactBuilder,
		stackService,
		happyConfig,
		awsBackend,
	)
}
func createStack(
	ctx context.Context,
	cmd *cobra.Command,
	stackName string,
	tags map[string]string,
	forceFlag bool,
	artifactBuilder ab.ArtifactBuilderIface,
	stackService *stackservice.StackService,
	happyConfig *config.HappyConfig,
	awsBackend *backend.Backend,
) error {
	// 1.) if the stack already exists and force flag is used, call the update function instead
	_, err := stackService.GetStack(ctx, stackName)
	if err == nil && forceFlag {
		logrus.Debugf("stack '%s' already exists, it will be updated", stackName)
		return runUpdate(cmd, []string{})
	}

	// 2.) add an entry to the stacklist
	stack, err := stackService.Add(ctx, stackName, dryRun)
	if err != nil {
		return errors.Wrap(err, "failed to add an entry to stacklist")
	}

	return updateStack(
		ctx,
		cmd,
		stackName,
		map[string]string{},
		force,
		artifactBuilder,
		stackService,
		happyConfig,
		awsBackend,
	)
}

func updateStack(
	ctx context.Context,
	cmd *cobra.Command,
	stackName string,
	tags map[string]string,
	forceFlag bool,
	artifactBuilder ab.ArtifactBuilderIface,
	stackService *stackservice.StackService,
	happyConfig *config.HappyConfig,
	awsBackend *backend.Backend,
) error {
	// 1.) build the docker image locally and push the images
	if createTag {
		err := artifactBuilder.BuildAndPush(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to build and push")
		}
	}

	// 2.) update the workspace's meta variables
	// TODO: is this used? the only thing I think some old happy environments use is the priority?
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
		return errors.Wrap(err, "failed to update the stack meta")
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
	stack.PrintOutputs(ctx)
	return nil
}

type validation func() error

func validateImageExists(ctx context.Context, skipCheckTag bool, ab artifact_builder.ArtifactBuilderIface) validation {
	return func() error {
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

func validateExistingStack(ctx context.Context, stackService *stackservice.StackService, stackName string, force bool) validation {
	return func() error {
		if force {
			return nil
		}

		_, err := stackService.GetStack(ctx, stackName)
		if err != nil {
			return errors.Wrap(err, "unable to get stack")
		}

		return nil
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
