package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pushImages []string
var tag string
var extraTag string
var composeEnvFile string

func init() {
	rootCmd.AddCommand(pushCmd)
	config.ConfigureCmdWithBootstrapConfig(pushCmd)

	pushCmd.Flags().StringSliceVar(&pushImages, "images", []string{}, "List of images to push to registry.")
	pushCmd.Flags().StringVar(&tag, "tag", "", "Tag name for existing docker image. Leave empty to generate one automatically.")
	pushCmd.Flags().StringVar(&extraTag, "extra-tag", "", "Extra tags to apply and push to the docker repo.")
	pushCmd.Flags().StringVar(&composeEnvFile, "compose-env", "", "Environment file to pass to docker compose")
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push docker images",
	Long:  "Push docker images to ECR",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPushWithOptions(tag, pushImages, extraTag)
	},
}

func runPush(tag string) error {
	return runPushWithOptions(tag, []string{}, "")
}

func runPushWithOptions(tag string, images []string, extraTag string) error {
	bootstrapConfig, err := config.NewBootstrapConfig()
	if err != nil {
		return err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return err
	}

	if len(composeEnvFile) == 0 {
		composeEnvFile = happyConfig.DefaultComposeEnvFile()
	}

	buildConfig := artifact_builder.NewBuilderConfig(bootstrapConfig, composeEnvFile, happyConfig.GetDockerRepo())
	artifactBuilder := artifact_builder.NewArtifactBuilder(buildConfig, happyConfig)
	serviceRegistries := happyConfig.GetRdevServiceRegistries()

	// NOTE login before build in order for cache to work
	err = artifactBuilder.RegistryLogin(serviceRegistries)
	if err != nil {
		return err
	}

	servicesImage, err := buildConfig.GetBuildServicesImage()
	if err != nil {
		return errors.Wrap(err, "failed to get service image")
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
	if len(extraTag) > 0 {
		allTags = append(allTags, extraTag)
	}
	fmt.Println(allTags)

	err = artifactBuilder.Build()
	if err != nil {
		return errors.Wrap(err, "failed to push image")
	}
	fmt.Println("Build complete")

	return artifactBuilder.Push(serviceRegistries, servicesImage, allTags)
}
