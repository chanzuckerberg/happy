package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
	config.ConfigureCmdWithBootstrapConfig(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build docker images",
	Long:  "Build docker images using docker-compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		bootstrapConfig, err := config.NewBootstrapConfig()
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(bootstrapConfig)
		if err != nil {
			return err
		}

		composeEnv := ""
		builderConfig := artifact_builder.NewBuilderConfig(bootstrapConfig, composeEnv)
		artifactBuilder := artifact_builder.NewArtifactBuilder(builderConfig, happyConfig)
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
