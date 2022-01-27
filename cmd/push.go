package cmd

import (
	"errors"
	"fmt"
	"os"

	// "time"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	// "github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"

	"github.com/spf13/cobra"
)

var pushImages []string
var tag string
var extraTag string
var composeEnv string

// TODO add support for this flag
var useComposeEnv bool

func init() {
	pushCmd.Flags().StringSliceVar(&pushImages, "images", []string{}, "List of images to push to registry.")
	pushCmd.Flags().StringVar(&tag, "tag", "", "Tag name for existing docker image. Leave empty to generate one automatically.")
	pushCmd.Flags().StringVar(&tag, "extra-tag", "", "Extra tags to apply and push to the docker repo.")
	pushCmd.Flags().StringVar(&composeEnv, "compose-env", "", "Environment file to pass to docker compose.")
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push docker images",
	Long:  "Push docker images to ECR",
	RunE: func(cmd *cobra.Command, args []string) error {
		updateCmd.Flags().StringVar(&sliceDefaultTag, "slice-default-tag", "", "For stacks using slices, override the default tag for any images that aren't being built & pushed by the slice")
		return runPushWithOptions(tag, pushImages, extraTag, composeEnv)
	},
}

func runPush(tag string) error {
	return runPushWithOptions(tag, make([]string, 0), "", "")
}

func runPushWithOptions(tag string, images []string, extraTag string, composeEnv string) error {
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
		return fmt.Errorf("failed to get Happy Config: %s", err)
	}

	if useComposeEnv {
		if len(composeEnv) == 0 {
			composeEnv = happyConfig.DefaultComposeEnv()
		}
	}

	buildConfig := artifact_builder.NewBuilderConfig(dockerComposeConfigPath, composeEnv)
	artifactBuilder := artifact_builder.NewArtifactBuilder(buildConfig, happyConfig)
	serviceRegistries, err := happyConfig.GetRdevServiceRegistries()
	if err != nil {
		return err
	}
	// NOTE login before build in order for cache to work
	artifactBuilder.RegistryLogin(serviceRegistries, pushImages)

	servicesImage, err := buildConfig.GetBuildServicesImage()
	if err != nil {
		return fmt.Errorf("failed to get service image: %s", err)
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
	allTags := []string{tag}
	fmt.Println(allTags)

	err = artifactBuilder.Build()
	if err != nil {
		return fmt.Errorf("failed to push image: %s", err)
	}
	fmt.Println("Build complete")

	// TODO add extra tag from input

	artifactBuilder.Push(serviceRegistries, servicesImage, allTags)
	return nil
}
