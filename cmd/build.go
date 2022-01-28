package cmd

import (
	"os"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build docker images",
	Long:  "Build docker images using docker-compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		dockerComposeConfig, ok := os.LookupEnv("DOCKER_COMPOSE_CONFIG_PATH")
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
		buildConfig := artifact_builder.NewBuilderConfig(dockerComposeConfig, composeEnv)
		artifactBuilder := artifact_builder.NewArtifactBuilder(buildConfig, happyConfig)
		serviceRegistries, err := happyConfig.GetRdevServiceRegistries()
		if err != nil {
			return err
		}
		// NOTE  not to login before build for cache to work
		err = artifactBuilder.RegistryLogin(serviceRegistries)
		if err != nil {
			return err
		}

		return artifactBuilder.Build()
	},
}
