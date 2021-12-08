package cmd

import (
	"errors"
	"os"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
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

		env := "rdev"

		dockerComposeConfig, ok := os.LookupEnv("DOCKER_COMPOSE_CONFIG_PATH")
		if !ok {
			return errors.New("Please set env var DOCKER_COMPOSE_CONFIG_PATH")
		}

		happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
		if !ok {
			return errors.New("Please set env var HAPPY_CONFIG_PATH")
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
		buildImages := []string{}
		artifactBuilder.RegistryLogin(serviceRegistries, buildImages)

		artifactBuilder.Build()
		return nil
	},
}
