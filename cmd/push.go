package cmd

import (
	"fmt"
	"os"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pushImages []string

// TODO add support for this flag
var useComposeEnv bool

func init() {
	pushCmd.Flags().StringSliceVar(&pushImages, "images", []string{}, "List of images to push to registry.")
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push docker images",
	Long:  "Push docker images to ECR",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: why is this empty?
		tag := ""
		return runPush(tag)
	},
}

func runPush(tag string) error {
	// TODO do not hardcode dev
	env := "rdev"

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
		return errors.Errorf("failed to get Happy Config: %s", err)
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

	// NOTE login before build in order for cache to work
	err = artifactBuilder.RegistryLogin(serviceRegistries)
	if err != nil {
		return err
	}

	servicesImage, err := buildConfig.GetBuildServicesImage()
	if err != nil {
		return errors.Errorf("failed to get service image: %s", err)
	}

	for service, reg := range serviceRegistries {
		fmt.Printf("%q: %q\t%q\n", service, reg.GetRepoUrl(), reg.GetRegistryUrl())
	}

	if tag == "" {
		tag, err = util.GenerateTag(happyConfig)
		if err != nil {
			return err
		}
	}
	tags := []string{tag}
	fmt.Println(tags)

	err = artifactBuilder.Build()
	if err != nil {
		return errors.Errorf("failed to push image: %s", err)
	}
	fmt.Println("Build complete")

	// TODO add extra tag from input
	return artifactBuilder.Push(serviceRegistries, servicesImage, tags)
}
