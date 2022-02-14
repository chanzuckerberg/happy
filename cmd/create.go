package cmd

import (
	"fmt"
	"strings"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/options"
	"github.com/chanzuckerberg/happy/pkg/orchestrator"
	stack_service "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	force        bool
	skipCheckTag bool
	createTag    bool
	tag          string
)

func init() {
	rootCmd.AddCommand(createCmd)
	config.ConfigureCmdWithBootstrapConfig(createCmd)

	createCmd.Flags().StringVar(&tag, "tag", "", "Specify the tag for the docker images. If not specified we will generate a default tag.")
	createCmd.Flags().BoolVar(&createTag, "create-tag", true, "Will build, tag, and push images when set. Otherwise, assumes images already exist.")
	createCmd.Flags().BoolVar(&force, "force", false, "Ignore the already-exists errors")
	createCmd.Flags().BoolVar(&skipCheckTag, "skip-check-tag", false, "Skip checking that the specified tag exists (requires --tag)")
}

var createCmd = &cobra.Command{
	Use:          "create STACK_NAME",
	Short:        "create new stack",
	Long:         "Create a new stack with a given tag.",
	SilenceUsage: true,
	PreRunE:      checkCreateFlags,
	RunE:         runCreate,
	Args:         cobra.ExactArgs(1),
}

func checkCreateFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("skip-check-tag") && !cmd.Flags().Changed("tag") {
		return errors.New("--skip-check-tag can only be used when --tag is specified")
	}

	if !createTag && !cmd.Flags().Changed("tag") {
		return errors.New("Must specify a tag when create-tag=false")
	}

	return nil
}

func runCreate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	stackName := args[0]

	bootstrapConfig, err := config.NewBootstrapConfig()
	if err != nil {
		return err
	}
	happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
	if err != nil {
		return err
	}

	backend, err := backend.NewAWSBackend(ctx, happyConfig)
	if err != nil {
		return err
	}

	builderConfig := artifact_builder.NewBuilderConfig(bootstrapConfig, happyConfig)
	ab := artifact_builder.NewArtifactBuilder(builderConfig, backend)

	url := backend.Conf().GetTfeUrl()
	org := backend.Conf().GetTfeOrg()

	workspaceRepo, err := workspace_repo.NewWorkspaceRepo(url, org)
	if err != nil {
		return err
	}
	stackService := stack_service.NewStackService(backend, workspaceRepo)

	existingStacks, err := stackService.GetStacks(ctx)
	if err != nil {
		return err
	}
	if _, ok := existingStacks[stackName]; ok {
		if !force {
			return errors.Errorf("stack %s already exists", stackName)
		}
	}

	// if creating tag and none specified, generate the default tag
	if createTag && len(tag) == 0 {
		tag, err = backend.GenerateTag(ctx)
		if err != nil {
			return err
		}
	}

	// if creating tag, build and push images
	if createTag {
		err = ab.BuildAndPush(ctx, artifact_builder.WithTags(tag))
		if err != nil {
			return err
		}
	}

	// check if image exists unless asked not to
	if !skipCheckTag {
		exists, err := ab.CheckImageExists(tag)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Errorf("image tag does not exist or cannot be verified: %s", createTag)
		}
	}

	// now that we have images, create all TFE related resources
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

	err = stackMeta.Update(ctx, tag, map[string]string{}, "", stackService)
	if err != nil {
		return err
	}
	fmt.Printf("Creating %s\n", stackName)

	stack, err := stackService.Add(ctx, stackName)
	if err != nil {
		return err
	}
	fmt.Printf("setting stackMeta %v\n", stackMeta)
	stack.SetMeta(stackMeta)

	err = stack.Apply(getWaitOptions(backend, stackName))
	if err != nil {
		return err
	}

	autoRunMigration := happyConfig.AutoRunMigration()
	if autoRunMigration {
		err = runMigrate(ctx, stackName)
		if err != nil {
			return err
		}
	}

	stack.PrintOutputs()
	return nil
}

// func buildSlice(
// 	ctx context.Context,
// 	backend *backend.Backend,
// 	sliceName string,
// 	defaultSliceTag string,
// ) (stackTags map[string]string, defaultTag string, err error) {
// 	defaultTag = defaultSliceTag

// 	slice, err := backend.Conf().GetSlice(sliceName)
// 	if err != nil {
// 		return nil, defaultTag, err
// 	}

// 	buildProfile := slice.Profile
// 	sliceTag, err := backend.GenerateTag(ctx)
// 	if err != nil {
// 		return stackTags, defaultTag, err
// 	}

// 	err = runPushWithOptions(ctx, sliceTag, buildProfile, "")
// 	if err != nil {
// 		return stackTags, defaultTag, errors.Wrap(err, "failed to push image")
// 	}

// 	if len(defaultTag) == 0 {
// 		defaultTag = backend.Conf().SliceDefaultTag()
// 	}

// 	stackTags = make(map[string]string)
// 	for _, sliceImg := range buildImages {
// 		stackTags[sliceImg] = sliceTag
// 	}

// 	return stackTags, defaultTag, nil
// }

func joinKeys(m map[string]config.Slice, separator string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, separator)
}

func getWaitOptions(backend *backend.Backend, stackName string) options.WaitOptions {
	taskOrchestrator := orchestrator.NewOrchestrator(backend)
	waitOptions := options.WaitOptions{StackName: stackName, Orchestrator: taskOrchestrator, Services: backend.Conf().GetServices()}
	return waitOptions
}
