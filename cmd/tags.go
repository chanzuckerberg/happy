package cmd

import (
	"os"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var sourceTag string
var destTags []string
var images []string

func init() {
	tagsCmd.Flags().StringVar(&sourceTag, "source-tag", "", "Tag name for existing docker image.")
	tagsCmd.Flags().StringSliceVar(&destTags, "dest-tag", []string{}, "Extra tags to apply and push to the docker repo.")
	rootCmd.AddCommand(tagsCmd)
}

var tagsCmd = &cobra.Command{
	Use:   "addtags",
	Short: "Add additional tags to already-pushed images in the ECR repo",
	Long:  "Add additional tags to already-pushed images in the ECR repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		env := "rdev"
		images = args

		dockerComposeConfigPath, ok := os.LookupEnv("DOCKER_COMPOSE_CONFIG_PATH")
		if !ok {
			return errors.New("please set env var DOCKER_COMPOSE_CONFIG_PATH")
		}

		happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
		if !ok {
			return errors.New("please set env var HAPPY_CONFIG_PATH")
		}

		happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
		if err != nil {
			return err
		}

		composeEnv := ""
		if useComposeEnv {
			composeEnv = happyConfig.DefaultComposeEnv()
		}

		buildConfig := artifact_builder.NewBuilderConfig(dockerComposeConfigPath, composeEnv)
		artifactBuilder := artifact_builder.NewArtifactBuilder(buildConfig, happyConfig)
		serviceRegistries, err := happyConfig.GetRdevServiceRegistries()
		if err != nil {
			return err
		}

		servicesImage, err := buildConfig.GetBuildServicesImage()
		if err != nil {
			return errors.Errorf("failed to get service image: %s", err)
		}

		return artifactBuilder.RetagImages(serviceRegistries, servicesImage, sourceTag, destTags, images)
	},
}
