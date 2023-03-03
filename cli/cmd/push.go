package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/spf13/cobra"
)

var tags []string
var sliceName string

func init() {
	rootCmd.AddCommand(pushCmd)
	config.ConfigureCmdWithBootstrapConfig(pushCmd)

	pushCmd.Flags().StringVar(&sliceName, "slice", "", "The name of the slice you'd like to push to the registry.")
	pushCmd.Flags().StringSliceVar(&tags, "tag", nil, "Extra tags to set for built images, comma-delimited (ex: tag1,tag2,tag3). We will, in addition, generate default tags automatically.")
}

var pushCmd = &cobra.Command{
	Use:          "push",
	Short:        "Push docker images",
	Long:         "Push docker images to ECR",
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

		b, err := backend.NewAWSBackend(ctx, happyConfig)
		if err != nil {
			return err
		}

		buildConfig := artifact_builder.NewBuilderConfig().WithBootstrap(bootstrapConfig).WithHappyConfig(happyConfig)
		// FIXME: this is an error-prone interface
		if sliceName != "" {
			slice, err := happyConfig.GetSlice(sliceName)
			if err != nil {
				return err
			}
			buildConfig.Profile = slice.Profile
		}

		artifactBuilder := artifact_builder.CreateArtifactBuilder().WithConfig(buildConfig).WithBackend(b).WithTags(tags)

		return artifactBuilder.BuildAndPush(ctx)
	},
}
