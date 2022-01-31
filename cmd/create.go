package cmd

import (
	"fmt"
	"os"
	"strings"

	"log"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	stack_service "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	createTag       string
	wait            bool
	force           bool
	sliceName       string
	sliceDefaultTag string
	skipCheckTag    bool
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createTag, "tag", "", "Tag name for docker image. Leave empty to generate one")
	createCmd.Flags().BoolVar(&wait, "wait", true, "Wait for this cmd to complete")
	createCmd.Flags().BoolVar(&force, "force", false, "Ignore the already-exists errors")
	createCmd.Flags().StringVarP(&sliceName, "slice", "s", "", "If you only need to test a slice of the app, specify it here")
	createCmd.Flags().StringVar(&sliceDefaultTag, "slice-default-tag", "", "For stacks using slices, override the default tag for any images that aren't being built & pushed by the slice")
	createCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
}

var createCmd = &cobra.Command{
	Use:   "create STACK_NAME",
	Short: "create new stack",
	Long:  "Create a new stack with a given tag.",
	RunE:  runCreate,
	Args:  cobra.ExactArgs(1),
}

func runCreate(cmd *cobra.Command, args []string) error {
	env := "rdev"

	stackName := args[0]

	fmt.Printf("Creating %s with settings: wait=%v force=%v\n", stackName, wait, force)

	dockerComposeConfigPath, ok := os.LookupEnv("DOCKER_COMPOSE_CONFIG_PATH")
	if !ok {
		return errors.New("please set env var DOCKER_COMPOSE_CONFIG_PATH")
	}

	happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
	if !ok {
		return errors.New("please set env var HAPPY_CONFIG_PATH")
	}

	_, ok = os.LookupEnv("HAPPY_PROJECT_ROOT")
	if !ok {
		return errors.New("please set env var HAPPY_PROJECT_ROOT")
	}

	happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
	if err != nil {
		return err
	}

	url, err := happyConfig.TfeUrl()
	if err != nil {
		return err
	}
	org, err := happyConfig.TfeOrg()
	if err != nil {
		return err
	}
	workspaceRepo, err := workspace_repo.NewWorkspaceRepo(url, org)
	if err != nil {
		return err
	}
	paramStoreBackend := backend.GetAwsBackend(happyConfig)
	stackService := stack_service.NewStackService(happyConfig, paramStoreBackend, workspaceRepo)

	existingStacks, err := stackService.GetStacks()
	if err != nil {
		return err
	}
	if _, ok := existingStacks[stackName]; ok {
		if !force {
			return errors.Errorf("stack %s already exists", stackName)
		}
	}

	stackTags := map[string]string{}
	if len(sliceName) > 0 {
		stackTags, createTag, err = buildSlice(happyConfig, sliceName, sliceDefaultTag)
		if err != nil {
			return err
		}
	}

	if !checkImageExists(dockerComposeConfigPath, env, happyConfig, tag) {
		return errors.Errorf("image tag does not exist or cannot be verified: %s", tag)
	}

	stackMeta := stackService.NewStackMeta(stackName)
	secretArn := happyConfig.GetSecretArn()
	if err != nil {
		return err
	}
	metaTag := map[string]string{"happy/meta/configsecret": secretArn}
	err = stackMeta.Load(metaTag)
	if err != nil {
		return err
	}

	if createTag == "" {
		createTag, err = util.GenerateTag(happyConfig)
		if err != nil {
			return err
		}

		// invoke push cmd
		fmt.Printf("Pushing images with tags %s...\n", createTag)
		err := runPush(createTag)
		if err != nil {
			return errors.Errorf("failed to push image: %s", err)
		}
	}
	err = stackMeta.Update(createTag, stackTags, sliceName, stackService)
	if err != nil {
		return err
	}
	fmt.Printf("Creating %s\n", stackName)

	stack, err := stackService.Add(stackName)
	if err != nil {
		return err
	}
	fmt.Printf("setting stackMeta %v\n", stackMeta)
	stack.SetMeta(stackMeta)

	err = stack.Apply()
	if err != nil {
		return err
	}

	autoRunMigration := happyConfig.AutoRunMigration()
	if err != nil {
		fmt.Println("WARNING autoRunMigration flag not set, defaulting to false")
	}

	if autoRunMigration {
		err = runMigrate(stackName)
		if err != nil {
			return err
		}
	}

	// TODO migrate db here
	stack.PrintOutputs()
	return nil
}

func checkImageExists(composeFile string, env string, happyConfig config.HappyConfig, tag string) bool {
	if len(tag) == 0 && skipCheckTag {
		return true
	}
	// Make sure all of our service images actually exist if we're trying to deploy via tag

	buildConfig := artifact_builder.NewBuilderConfig(composeFile, env)
	artifactBuilder := artifact_builder.NewArtifactBuilder(buildConfig, happyConfig)

	serviceRegistries, err := happyConfig.GetRdevServiceRegistries()
	if err != nil {
		log.Printf("Unable to retrieve service container registry information: %s\n", err.Error())
		return false
	}

	return artifactBuilder.CheckImageExists(serviceRegistries, tag)
}

func buildSlice(happyConfig config.HappyConfig, sliceName string, defaultSliceTag string) (stackTags map[string]string, defaultTag string, err error) {
	defaultTag = defaultSliceTag

	slices, err := happyConfig.GetSlices()
	if err != nil {
		return stackTags, defaultTag, errors.Errorf("unable to retrieve slice configuration: %s", err.Error())
	}

	slice, ok := slices[sliceName]
	if !ok {
		validSlices := joinKeys(slices, ", ")
		return stackTags, defaultTag, errors.Errorf("slice %s is invalid - valid names: %s", sliceName, validSlices)
	}

	buildImages := slice.BuildImages
	sliceTag, err := util.GenerateTag(happyConfig)
	if err != nil {
		return stackTags, defaultTag, err
	}

	err = runPushWithOptions(sliceTag, buildImages, "", "")
	if err != nil {
		return stackTags, defaultTag, errors.Errorf("failed to push image: %s", err.Error())
	}

	if len(defaultTag) == 0 {
		defaultTag = happyConfig.SliceDefaultTag()
	}

	stackTags = make(map[string]string)
	for _, sliceImg := range buildImages {
		stackTags[sliceImg] = sliceTag
	}

	return stackTags, defaultTag, nil
}

func joinKeys(m map[string]config.Slice, separator string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, separator)
}
