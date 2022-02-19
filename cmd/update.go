package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sliceDefaultTag string

func init() {
	rootCmd.AddCommand(updateCmd)
	config.ConfigureCmdWithBootstrapConfig(updateCmd)

	updateCmd.Flags().StringVar(&tag, "tag", "", "Tag name for docker image. Leave empty to generate one automatically.")
	updateCmd.Flags().StringVarP(&sliceName, "slice", "s", "", "If you only need to test a slice of the app, specify it here")
	updateCmd.Flags().StringVar(&sliceDefaultTag, "slice-default-tag", "", "For stacks using slices, override the default tag for any images that aren't being built & pushed by the slice")
	updateCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
}

var updateCmd = &cobra.Command{
	Use:          "update STACK_NAME",
	Short:        "update stack",
	Long:         "Update stack mathcing STACK_NAME",
	SilenceUsage: true,
	// PreRunE:      checkFlags,
	RunE: runUpdate,
	Args: cobra.ExactArgs(1),
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

	b, err := backend.NewAWSBackend(ctx, happyConfig)
	if err != nil {
		return err
	}

	builderConfig := artifact_builder.NewBuilderConfig().WithBootstrap(bootstrapConfig).WithHappyConfig(happyConfig)
	ab := artifact_builder.NewArtifactBuilder().WithConfig(builderConfig).WithBackend(b)

	url := b.Conf().GetTfeUrl()
	org := b.Conf().GetTfeOrg()

	workspaceRepo, err := workspace_repo.NewWorkspaceRepo(url, org)
	if err != nil {
		return err
	}

	stackService := stackservice.NewStackService(b, workspaceRepo)

	// build and push; creating tag if needed
	if tag == "" {
		tag, err = b.GenerateTag(ctx)
		if err != nil {
			return err
		}
	}

	buildOpts := []artifact_builder.ArtifactBuilderBuildOption{
		artifact_builder.WithTags(tag),
	}

	// if slice specified, use it
	if sliceName != "" {
		slice, err := happyConfig.GetSlice(sliceName)
		if err != nil {
			return err
		}
		buildOpts = append(buildOpts, artifact_builder.BuildSlice(slice))
	}
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

		// TODO: 2x check this
		for _, image := range serviceImages {
			stackTags[image] = tag
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

	logrus.Infof("updating stack %s", stackName)

	stacks, err := stackService.GetStacks(ctx)
	if err != nil {
		return err
	}
	stack, ok := stacks[stackName]
	if !ok {
		return errors.Errorf("stack %s not found", stackName)
	}

	stackMeta, err := stack.Meta()
	if err != nil {
		return err
	}

	// reset the configsecret if it has changed
	secretArn := happyConfig.GetSecretArn()

	configSecret := map[string]string{"happy/meta/configsecret": secretArn}
	err = stackMeta.Load(configSecret)
	if err != nil {
		return err
	}

	err = stackMeta.Update(ctx, tag, stackTags, sliceDefaultTag, stackService)
	if err != nil {
		return err
	}

	err = stack.Apply(getWaitOptions(b, stackName))
	if err != nil {
		return errors.Wrap(err, "apply failed, skipping migrations")
	}

	stack.PrintOutputs()
	return nil
}
