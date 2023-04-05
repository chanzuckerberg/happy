package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
	config.ConfigureCmdWithBootstrapConfig(buildCmd)
	cmd.SupportBuildSlices(buildCmd, &sliceName, &sliceDefaultTag)
}

var buildCmd = &cobra.Command{
	Use:          "build",
	Short:        "Build docker images",
	Long:         "Build docker images using docker compose",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		bootstrapConfig, err := config.NewBootstrapConfig(cmd)
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(bootstrapConfig)
		if err != nil {
			return err
		}
		backend, err := aws.NewAWSBackend(ctx, happyConfig.GetEnvironmentContext())
		if err != nil {
			return err
		}

		builderConfig := artifact_builder.NewBuilderConfig().
			WithBootstrap(bootstrapConfig).
			WithHappyConfig(happyConfig)
		// FIXME: this is an error-prone interface
		if sliceName != "" {
			slice, err := happyConfig.GetSlice(sliceName)
			if err != nil {
				return err
			}
			builderConfig.Profile = slice.Profile
		}
		artifactBuilder := artifact_builder.CreateArtifactBuilder().
			WithHappyConfig(happyConfig).
			WithConfig(builderConfig).
			WithBackend(backend)
		// NOTE not to login before build for cache to work
		err = artifactBuilder.RegistryLogin(ctx)
		if err != nil {
			return err
		}

		return artifactBuilder.Build(ctx)
	},
}
