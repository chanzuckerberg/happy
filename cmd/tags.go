package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
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
}

var tagsCmd = &cobra.Command{
	Use:   "addtags",
	Short: "Add additional tags to already-pushed images in the ECR repo",
	Long:  "Add additional tags to already-pushed images in the ECR repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		images = args

		bootstrapConfig, err := config.NewBootstrapConfig()
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(bootstrapConfig)
		if err != nil {
			return err
		}

		composeEnv := happyConfig.DefaultComposeEnv()
		buildConfig := artifact_builder.NewBuilderConfig(bootstrapConfig, composeEnv, happyConfig.GetDockerRepo())
		artifactBuilder := artifact_builder.NewArtifactBuilder(buildConfig, happyConfig)
		serviceRegistries, err := happyConfig.GetRdevServiceRegistries()
		if err != nil {
			return err
		}

		servicesImage, err := buildConfig.GetBuildServicesImage()
		if err != nil {
			return errors.Wrap(err, "failed to get service image")
		}

		return artifactBuilder.RetagImages(serviceRegistries, servicesImage, sourceTag, destTags, images)
	},
}
