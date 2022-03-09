package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/spf13/cobra"
)

var sourceTag string
var destTags []string
var images []string

func init() {
	rootCmd.AddCommand(tagsCmd)
	config.ConfigureCmdWithBootstrapConfig(tagsCmd)

	tagsCmd.Flags().StringVar(&sourceTag, "source-tag", "", "Tag name for existing docker image.")
	tagsCmd.Flags().StringSliceVar(&destTags, "dest-tag", []string{}, "Extra tags to apply and push to the docker repo.")
	_ = tagsCmd.MarkFlagRequired("source-tag")
	_ = tagsCmd.MarkFlagRequired("dest-tag")
}

var tagsCmd = &cobra.Command{
	Use:          "addtags",
	Short:        "Add additional tags to already-pushed images in the ECR repo",
	Long:         "Add additional tags to already-pushed images in the ECR repo",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		images = args

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
		artifactBuilder := artifact_builder.NewArtifactBuilder().WithConfig(buildConfig).WithBackend(b)
		serviceRegistries := b.Conf().GetServiceRegistries()

		return artifactBuilder.RetagImages(serviceRegistries, sourceTag, destTags, images)
	},
}
