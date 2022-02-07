package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pushImages []string
var tag string
var extraTag string

func init() {
	rootCmd.AddCommand(pushCmd)
	config.ConfigureCmdWithBootstrapConfig(pushCmd)

	pushCmd.Flags().StringSliceVar(&pushImages, "images", []string{}, "List of images to push to registry.")
	pushCmd.Flags().StringVar(&tag, "tag", "", "Tag name for existing docker image. Leave empty to generate one automatically.")
	pushCmd.Flags().StringVar(&extraTag, "extra-tag", "", "Extra tags to apply and push to the docker repo.")
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

	composeEnv := happyConfig.DefaultComposeEnv()
	buildConfig := artifact_builder.NewBuilderConfig(bootstrapConfig, composeEnv, happyConfig.GetDockerRepo())
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
		return errors.Wrap(err, "failed to get service image")
	}

	for service, reg := range serviceRegistries {
		log.Printf("%q: %q\t%q\n", service, reg.GetRepoUrl(), reg.GetRegistryUrl())
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
	log.Println(allTags)

	err = artifactBuilder.Build()
	if err != nil {
		return errors.Wrap(err, "failed to push image")
	}
	log.Println("Build complete")

	return artifactBuilder.Push(serviceRegistries, servicesImage, allTags)
}
