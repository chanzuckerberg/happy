package cmd

import (
	"context"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	happyCmd "github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	waitoptions "github.com/chanzuckerberg/happy/pkg/options"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/hashicorp/go-multierror"
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
	createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Prepare all resources, but do not apply any changes")
}

var createCmd = &cobra.Command{
	Use:          "create STACK_NAME",
	Short:        "create new stack",
	Long:         "Create a new stack with a given tag.",
	SilenceUsage: true,
	PreRunE:      happyCmd.Validate(checkCreateFlags, cobra.ExactArgs(1), happyCmd.CheckStackName),
	RunE:         runCreate,
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

func runCreate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	stackName := args[0]

	if !regexp.MustCompile(`^[a-z0-9\-]*$`).MatchString(stackName) {
		return errors.New("Stack name must contain only lowercase letters, numbers, and hyphens/dashes")
	}

	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
	if err != nil {
		return err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return err
	}

	backend, err := backend.NewAWSBackend(ctx, happyConfig)
	if err != nil {
		return err
	}

	builderConfig := artifact_builder.NewBuilderConfig().WithBootstrap(bootstrapConfig).WithHappyConfig(happyConfig)

	// slice support parity with update command
	buildOpts := []artifact_builder.ArtifactBuilderBuildOption{}
	// FIXME: this is an error-prone interface
	// if slice specified, use it
	if sliceName != "" {
		slice, err := happyConfig.GetSlice(sliceName)
		if err != nil {
			return err
		}
		buildOpts = append(buildOpts, artifact_builder.BuildSlice(slice))
		builderConfig.WithProfile(slice.Profile)
	}
	ab := artifact_builder.NewArtifactBuilder().WithConfig(builderConfig).WithBackend(backend)

	url := backend.Conf().GetTfeUrl()
	org := backend.Conf().GetTfeOrg()

	workspaceRepo := workspace_repo.NewWorkspaceRepo(url, org).WithDryRun(dryRun)
	stackService := stackservice.NewStackService().WithBackend(backend).WithWorkspaceRepo(workspaceRepo)

	err = verifyTFEBacklog(ctx, workspaceRepo)
	if err != nil {
		return err
	}

	existingStacks, err := stackService.GetStacks(ctx)
	if err != nil {
		return err
	}
	existingStack, stackAlreadyExists := existingStacks[stackName]
	if stackAlreadyExists && !force {
		return errors.Errorf("stack '%s' already exists, use 'happy update %s' to update it", stackName, stackName)
	}

	if dryRun {
		createTag = false
		skipCheckTag = true
	}

	// if creating tag and none specified, generate the default tag
	if createTag && (tag == "") {
		tag, err = backend.GenerateTag(ctx)
		if err != nil {
			return err
		}
	}

	// if creating tag, build and push images
	if createTag {
		// parity with update command
		buildOpts = append(buildOpts, artifact_builder.WithTags(tag))
		err = ab.BuildAndPush(ctx, buildOpts...)
		if err != nil {
			return err
		}
	}

	// check if image exists unless asked not to
	if !skipCheckTag {
		exists, err := ab.CheckImageExists(ctx, tag)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Errorf("image tag does not exist: '%s'", tag)
		}
	}

	// if we already have a stack and "force" then use existing
	var stackMeta *stackservice.StackMeta
	if force && existingStack != nil {
		stackMeta, err = existingStack.Meta(ctx)
		if err != nil {
			return err
		}
	} else {
		stackMeta = stackService.NewStackMeta(stackName)
	}

	options := stackservice.NewStackManagementOptions(stackName).WithHappyConfig(happyConfig).WithStackService(stackService).WithStackMeta(stackMeta).WithBackend(backend)

	// now that we have images, create all TFE related resources
	return createStack(ctx, cmd, options)
}

func createStack(ctx context.Context, cmd *cobra.Command, options *stackservice.StackManagementOptions) error {
	var errs *multierror.Error

	if options.Stack != nil {
		errs = multierror.Append(errs, errors.New("stack option not expected in this context"))
	}
	if options.StackService == nil {
		errs = multierror.Append(errs, errors.New("stackService option not provided"))
	}
	if options.Backend == nil {
		errs = multierror.Append(errs, errors.New("backend option not provided"))
	}
	if options.StackMeta == nil {
		errs = multierror.Append(errs, errors.New("stackMeta option not provided"))
	}
	if len(options.StackName) == 0 {
		errs = multierror.Append(errs, errors.New("stackName option not provided"))
	}

	err := errs.ErrorOrNil()
	if err != nil {
		return err
	}

	secretArn := options.HappyConfig.GetSecretArn()

	metaTag := map[string]string{"happy/meta/configsecret": secretArn}
	err = options.StackMeta.Load(metaTag)
	if err != nil {
		return errors.Wrap(err, "failed to load stack meta")
	}

	targetBaseTag := tag
	if sliceDefaultTag != "" {
		targetBaseTag = sliceDefaultTag
	}

	err = options.StackMeta.Update(ctx, targetBaseTag, options.StackTags, "", options.StackService)
	if err != nil {
		return errors.Wrap(err, "failed to update the stack meta")
	}

	if dryRun {
		logrus.Infof("temporarily creating a TFE workspace for stack '%s'", options.StackName)
	} else {
		logrus.Infof("creating stack '%s'", options.StackName)
	}
	stack, err := options.StackService.Add(ctx, options.StackName)
	if err != nil {
		return errors.Wrap(err, "failed to add the stack")
	}

	logrus.Debugf("setting stackMeta %v", options.StackMeta)
	stack = stack.WithMeta(options.StackMeta)

	err = stack.Apply(ctx, getWaitOptions(options, dryRun))
	if err != nil {
		return errors.Wrap(err, "failed to successfully create the stack")
	}

	if dryRun {
		logrus.Infof("cleaning up stack '%s'", options.StackName)
		err = options.StackService.Remove(ctx, options.StackName)
		if err != nil {
			logrus.Errorf("failed to clean up the stack: %s", err.Error())
		}
	} else {
		shouldRunMigration, err := happyCmd.ShouldRunMigrations(cmd, options.HappyConfig)
		if err != nil {
			return err
		}
		if shouldRunMigration {
			err = runMigrate(cmd, options.StackName)
			if err != nil {
				return errors.Wrap(err, "failed to run migrations")
			}
		}
		stack.PrintOutputs(ctx)
	}

	return nil
}

func getWaitOptions(options *stackservice.StackManagementOptions, dryRun bool) waitoptions.WaitOptions {
	taskOrchestrator := orchestrator.NewOrchestrator().WithBackend(options.Backend)
	waitOptions := waitoptions.WaitOptions{
		StackName:    options.StackName,
		Orchestrator: taskOrchestrator,
		Services:     options.Backend.Conf().GetServices(),
		DryRun:       dryRun,
	}
	return waitOptions
}

func verifyTFEBacklog(ctx context.Context, workspaceRepo *workspace_repo.WorkspaceRepo) error {
	backlogSize, _, err := workspaceRepo.EstimateBacklogSize(ctx)
	if err != nil {
		return errors.Wrap(err, "error estimating TFE backlog")
	}
	if backlogSize < 2 {
		logrus.Info("There is no TFE backlog, proceeding.")
	} else if backlogSize < 20 {
		logrus.Infof("TFE backlog is only %d runs long, proceeding.", backlogSize)
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
