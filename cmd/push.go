package cmd

import (
	"context"
	"fmt"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
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
		return runPushWithOptions(cmd.Context(), tag, pushImages, extraTag)
	},
}

func runPush(ctx context.Context, tag string) error {
	return runPushWithOptions(ctx, tag, []string{}, "")
}

func runPushWithOptions(ctx context.Context, tag string, images []string, extraTag string) error {
	bootstrapConfig, err := config.NewBootstrapConfig()
	if err != nil {
		return err
	}
	happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
	if err != nil {
		return err
	}

	b, err := backend.NewAWSBackend(ctx, happyConfig)
	if err != nil {
		return err
	}

	buildConfig := artifact_builder.NewBuilderConfig(bootstrapConfig, happyConfig)
	artifactBuilder := artifact_builder.NewArtifactBuilder(buildConfig, b)
	serviceRegistries := b.Conf().GetServiceRegistries()

	// NOTE login before build in order for cache to work
	err = artifactBuilder.RegistryLogin(ctx, serviceRegistries)
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
		tag, err = b.GenerateTag(ctx)
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
