package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sliceDefaultTag string

func init() {
	rootCmd.AddCommand(updateCmd)
	config.ConfigureCmdWithBootstrapConfig(updateCmd)
	cmd.SupportUpdateSlices(updateCmd, &sliceName, &sliceDefaultTag)

	updateCmd.Flags().StringVar(&tag, "tag", "", "Tag name for docker image. Leave empty to generate one automatically.")
	updateCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
	updateCmd.Flags().BoolVar(&force, "force", false, "Force stack creation if it doesn't exist")
}

var updateCmd = &cobra.Command{
	Use:          "update STACK_NAME",
	Short:        "update stack",
	Long:         "Update stack mathcing STACK_NAME",
	SilenceUsage: true,
	RunE:         runUpdate,
	PreRunE:      cmd.Validate(cmd.ValidateUpdateSliceFlags, cobra.ExactArgs(1)),
}

func runUpdate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	stackName := args[0]

	bootstrapConfig, err := config.NewBootstrapConfig()
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

	ab := artifact_builder.NewArtifactBuilder().WithBackend(backend).WithConfig(builderConfig)
	defer ab.Profiler.PrintRuntimes()
	url := backend.Conf().GetTfeUrl()
	org := backend.Conf().GetTfeOrg()

	workspaceRepo := workspace_repo.NewWorkspaceRepo(url, org)
	stackService := stackservice.NewStackService().WithBackend(backend).WithWorkspaceRepo(workspaceRepo)

	// build and push; creating tag if needed
	if tag == "" {
		tag, err = backend.GenerateTag(ctx)
		if err != nil {
			return err
		}
	}

	buildOpts = append(buildOpts, artifact_builder.WithTags(tag))
	err = ab.BuildAndPush(ctx, buildOpts...)
	if err != nil {
		return err
	}

	// consolidate some stack tags
	stackTags := map[string]string{}
	if sliceName != "" {
		serviceImages, err := builderConfig.GetBuildServicesImage()
		if err != nil {
			return err
		}

		for service := range serviceImages {
			stackTags[service] = tag
		}
	}

	// check if image exists unless asked not to
	if !skipCheckTag {
		exists, err := ab.CheckImageExists(tag)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Errorf("image tag does not exist or cannot be verified: %s", tag)
		}
	}

	stacks, err := stackService.GetStacks(ctx)
	if err != nil {
		return err
	}

	options := stackservice.NewStackManagementOptions(stackName).WithHappyConfig(happyConfig).WithStackService(stackService).WithBackend(backend).WithStackTags(stackTags)

	stack, ok := stacks[stackName]
	if !ok {
		// Stack does not exist
		if !force {
			return errors.Errorf("stack '%s' does not exist, use --force or 'happy create %s' to create it", stackName, stackName)
		} else { // Force creation of the new stack
			logrus.Infof("stack '%s' doesn't exist, it will be created", stackName)
			stackMeta := stackService.NewStackMeta(stackName)
			options = options.WithStackMeta(stackMeta)
			return createStack(ctx, options)
		}
	}

	logrus.Infof("updating stack '%s'", stackName)
	options = options.WithStack(stack)

	// reset the configsecret if it has changed
	// if we have a default tag, use it
	err = updateStack(ctx, options)

	return err
}

func updateStack(ctx context.Context, options *stackservice.StackManagementOptions) error {
	var errs *multierror.Error

	if options.Stack == nil {
		errs = multierror.Append(errs, errors.New("stack option not provided"))
	}
	if options.StackService == nil {
		errs = multierror.Append(errs, errors.New("stackService option not provided"))
	}
	if options.Backend == nil {
		errs = multierror.Append(errs, errors.New("backend option not provided"))
	}
	if options.StackMeta != nil {
		errs = multierror.Append(errs, errors.New("stackMeta option should not be provided in this context"))
	}
	if len(options.StackName) == 0 {
		errs = multierror.Append(errs, errors.New("stackName option not provided"))
	}

	err := errs.ErrorOrNil()
	if err != nil {
		return err
	}

	stackMeta, err := options.Stack.Meta(ctx)
	if err != nil {
		return err
	}

	secretArn := options.HappyConfig.GetSecretArn()

	configSecret := map[string]string{"happy/meta/configsecret": secretArn}
	err = stackMeta.Load(configSecret)
	if err != nil {
		return err
	}

	targetBaseTag := tag
	if sliceDefaultTag != "" {
		targetBaseTag = sliceDefaultTag
	}

	err = stackMeta.Update(ctx, targetBaseTag, options.StackTags, sliceName, options.StackService)
	if err != nil {
		return err
	}

	err = options.Stack.Apply(ctx, getWaitOptions(options.Backend, options.StackName))
	if err != nil {
		return errors.Wrap(err, "apply failed, skipping migrations")
	}

	options.Stack.PrintOutputs(ctx)
	return nil
}
