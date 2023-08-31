package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/go-misc/sets"
	ab "github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	waitoptions "github.com/chanzuckerberg/happy/shared/options"
	stackservice "github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/util/tf"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type HappyClient struct {
	HappyConfig     *config.HappyConfig
	StackService    *stackservice.StackService
	ArtifactBuilder ab.ArtifactBuilderIface
	StackTags       map[string]string
	AWSBackend      *backend.Backend
}

func makeHappyClientFromBootstrap(ctx context.Context, bootstrapConfig *config.Bootstrap, sliceName, stackName string, tags []string, createTag bool) (*HappyClient, error) {
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, err
	}
	awsBackend, err := backend.NewAWSBackend(ctx, happyConfig.GetEnvironmentContext())
	if err != nil {
		return nil, err
	}
	builderConfig := ab.NewBuilderConfig().
		WithBootstrap(bootstrapConfig).
		WithHappyConfig(happyConfig)

	builderConfig.StackName = stackName
	ab, stackTags, err := configureArtifactBuilder(ctx, sliceName, tags, createTag, builderConfig, happyConfig, awsBackend)
	if err != nil {
		return nil, err
	}
	workspaceRepo := createWorkspaceRepo(awsBackend)
	stackService := stackservice.NewStackService(happyConfig.GetEnv(), happyConfig.App()).
		WithBackend(awsBackend).
		WithWorkspaceRepo(workspaceRepo)

	return &HappyClient{
		HappyConfig:     happyConfig,
		StackService:    stackService,
		ArtifactBuilder: ab,
		StackTags:       stackTags,
		AWSBackend:      awsBackend,
	}, nil
}

var happyClient *HappyClient

func makeHappyClient(cmd *cobra.Command, sliceName, stackName string, tags []string, createTag bool) (*HappyClient, error) {
	// reuse the happy client when possible so we don't call expensive auth operations multiple times
	if happyClient != nil {
		return happyClient, nil
	}
	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
	if err != nil {
		return nil, err
	}
	happyClient, err = makeHappyClientFromBootstrap(cmd.Context(), bootstrapConfig, sliceName, stackName, tags, createTag)
	return happyClient, err
}

func createWorkspaceRepo(backend *backend.Backend) workspace_repo.WorkspaceRepoIface {
	if util.IsLocalstackMode() {
		return workspace_repo.NewLocalWorkspaceRepo()
	}
	url := backend.Conf().GetTfeUrl()
	org := backend.Conf().GetTfeOrg()
	return workspace_repo.NewWorkspaceRepo(url, org)
}

func configureArtifactBuilder(
	ctx context.Context,
	sliceName string,
	tags []string,
	createTag bool,
	builderConfig *ab.BuilderConfig,
	happyConfig *config.HappyConfig,
	backend *backend.Backend) (ab.ArtifactBuilderIface, map[string]string, error) {
	artifactBuilder := ab.NewArtifactBuilder(ctx).
		WithHappyConfig(happyConfig).
		WithConfig(builderConfig).
		WithBackend(backend)
	var err error
	if sliceName != "" {
		slice, err := happyConfig.GetSlice(sliceName)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to find the slice %s", sliceName)
		}
		builderConfig.Profile = slice.Profile
	}

	// if creating tag and none specified, generate the default tag
	generatedTag := ""
	artifactBuilder.WithTags(tags)
	if createTag && len(artifactBuilder.GetTags()) == 0 {
		generatedTag, err = backend.GenerateTag(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to generate tag")
		}
		artifactBuilder.WithTags([]string{generatedTag})
	}

	stackTags := map[string]string{}
	if sliceName != "" {
		serviceImages, err := builderConfig.GetBuildServicesImage(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to get build service images")
		}

		for service := range serviceImages {
			stackTags[service] = generatedTag
		}
	}

	return artifactBuilder, stackTags, nil
}

type validation func() error

func validateImageExists(
	ctx context.Context,
	createTag, skipCheckTag bool,
	imageSrcEnv, imageSrcStack string,
	happyClient *HappyClient,
	useAWSProfile bool,
) validation {
	return func() error {
		logrus.Debug("Running validateImageExists()")
		if skipCheckTag {
			return nil
		}

		if imageSrcEnv != "" && imageSrcStack != "" {
			stackTag := strings.SplitN(imageSrcStack, ":", 2)
			tag := ""
			if len(stackTag) < 1 {
				return errors.Errorf("invalid image source stack %s", imageSrcStack)
			}

			stack := stackTag[0]
			if len(stackTag) > 1 {
				tag = stackTag[1]
			}

			// make a client associated with the env we are pulling from
			bs, err := config.NewBootstrapConfigForEnv(imageSrcEnv, useAWSProfile)
			if err != nil {
				return errors.Wrapf(err, "unable to bootstrap %s env", imageSrcEnv)
			}
			srcHappyClient, err := makeHappyClientFromBootstrap(ctx, bs, "", stack, []string{tag}, false)
			if err != nil {
				return errors.Wrapf(err, "unable to create happy client for env %s", imageSrcEnv)
			}

			return errors.Wrapf(pullAndPushImageFrom(ctx, stack, tag, srcHappyClient, happyClient),
				"unable to pull and push image from %s:%s",
				stack,
				tag)
		}

		if createTag {
			// if we build and push and it succeeds, we know that the image exists
			return happyClient.ArtifactBuilder.BuildAndPush(ctx)
		}

		if len(happyClient.ArtifactBuilder.GetTags()) == 0 {
			return errors.Errorf("no tags have been assigned")
		}

		for _, tag := range happyClient.ArtifactBuilder.GetTags() {
			exists, err := happyClient.ArtifactBuilder.CheckImageExists(ctx, tag)
			if err != nil {
				return errors.Wrapf(err, "error checking if tag %s existed", tag)
			}
			if !exists {
				return errors.Errorf("image tag does not exist: '%s'", tag)
			}
		}

		return nil
	}
}

func validateTFEBackLog(ctx context.Context, awsBackend *backend.Backend) validation {
	logrus.Debug("Scheduling validateTFEBackLog()")
	return func() error {
		logrus.Debug("Running validateTFEBackLog()")
		return verifyTFEBacklog(ctx, createWorkspaceRepo(awsBackend))
	}
}

func validateGitTree(projectRoot string) validation {
	logrus.Debug("Scheduling validateGitTree()")
	return func() error {
		logrus.Debug("Running validateGitTree()")
		return util.ValidateGitTree(projectRoot)
	}
}

func validateStackNameAvailable(ctx context.Context, stackService *stackservice.StackService, stackName string, force bool) validation {
	logrus.Debug("Scheduling validateStackNameAvailable()")
	return func() error {
		logrus.Debug("Running validateStackNameAvailable()")
		if force {
			return nil
		}

		metas, err := stackService.CollectStackInfo(ctx, happyClient.HappyConfig.App(), true)
		if err != nil {
			return errors.Wrap(err, "unable to collect stack info")
		}

		for _, meta := range metas {
			if meta.Stack == stackName {
				if meta.AppName == happyClient.HappyConfig.App() {
					return nil
				}
				return errors.Errorf("this stack exists, but in a different app ('%s'), you cannot manipulate it from this app", meta.AppName)
			}
		}

		return errors.Errorf("stack %s doesn't exist", stackName)
	}
}

func validateStackNameGloballyAvailable(ctx context.Context, stackService *stackservice.StackService, stackName string, force bool) validation {
	logrus.Debug("Scheduling validateStackNameAvailable()")
	return func() error {
		logrus.Debug("Running validateStackNameGloballyAvailable()")
		if force {
			return nil
		}

		metas, err := stackService.CollectStackInfo(ctx, happyClient.HappyConfig.App(), true)
		if err != nil {
			return errors.Wrap(err, "unable to collect stack info")
		}

		for _, meta := range metas {
			if meta.Stack == stackName {
				if meta.AppName == happyClient.HappyConfig.App() {
					return errors.New("the stack name is already taken by this app")
				}
				return errors.Errorf("the stack name is already taken by '%s' app; to see all stacks deployed, run 'happy list --all'", meta.AppName)
			}
		}

		return nil
	}
}

func validateConfigurationIntegirty(ctx context.Context, slice string, happyClient *HappyClient) validation {
	logrus.Debug("Scheduling validateConfigurationIntegirty()")
	return func() error {
		logrus.Debug("Running validateConfigurationIntegirty()")
		// These services are configured in docker-compose.yml, and have their containers built
		availableServices, err := happyClient.ArtifactBuilder.GetServices(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to get available services")
		}

		// NOTE: availableServices will only contain the services that are a part of the slice.
		// That means we have to iterate over those first to make sure they are all in the config.json and
		// not the other way around.
		configServices := happyClient.HappyConfig.GetServices()
		ss := sets.NewStringSet().Add(configServices...)
		for serviceName, service := range availableServices {
			// ignore services that don't have a build section
			// as these are not deployable services
			if service.Build == nil {
				continue
			}
			ok := ss.ContainsElement(serviceName)
			if !ok {
				return errors.Errorf("service %s was referenced in docker-compose.yml, but not referenced in .happy/config.json services array", serviceName)
			}
		}

		// These services are referenced in terraform code for the environment
		srcDir := filepath.Join(happyClient.HappyConfig.GetProjectRoot(), happyClient.HappyConfig.TerraformDirectory())
		deployedServices, err := tf.NewTfParser().ParseServices(srcDir)
		if err != nil {
			return errors.Wrap(err, "unable to parse terraform code")
		}
		for service := range deployedServices {
			if _, ok := availableServices[service]; !ok {
				return errors.Errorf("service %s is not configured in docker-compose.yml, but referenced in your terraform code", service)
			}
			found := false
			for _, configuredService := range happyClient.HappyConfig.GetServices() {
				if service == configuredService {
					found = true
					break
				}
			}

			if !found {
				return errors.Errorf("service %s is not configured in ./happy/config.json, but referenced in your terraform code", service)
			}
		}

		return nil
	}
}

// pull an image from a stack and push it to a new stack
// an optional tag can be provided to cherry pick a stack's image from the image history
func pullAndPushImageFrom(
	ctx context.Context,
	srcStackName, srcTag string, srcHappyClient *HappyClient,
	targetHappyClient *HappyClient,
) error {
	err := validate(validateStackExistsUpdate(ctx, srcStackName, srcHappyClient))
	if err != nil {
		return errors.Wrapf(err, "stack %s doesn't exist", srcStackName)
	}
	// if no tag is specified, get the latest deployed tag
	// by reading the image variable from the TFE workspace
	if srcTag == "" {
		var err error
		srcTag, err = srcHappyClient.StackService.GetLatestDeployedTag(ctx, srcStackName)
		if err != nil {
			return errors.Wrapf(err, "unable to get latest tag from stack %s", srcStackName)
		}
	} else {
		exists, err := srcHappyClient.ArtifactBuilder.CheckImageExists(ctx, srcTag)
		if err != nil {
			return errors.Wrapf(err, "error checking if tag %s existed", srcTag)
		}
		if !exists {
			return errors.Errorf("src image tag does not exist %s:%s", srcStackName, srcTag)
		}
	}

	servicesImage, err := srcHappyClient.ArtifactBuilder.Pull(ctx, srcStackName, srcTag)
	if err != nil {
		return errors.Wrapf(err, "unable to pull image %s from stack %s in env %s", srcTag, srcStackName, srcHappyClient.HappyConfig.GetEnv())
	}
	err = targetHappyClient.ArtifactBuilder.PushFromWithTag(ctx, servicesImage, srcTag)
	if err != nil {
		return errors.Wrapf(err, "unable to push image %s from stack %s in env %s", srcTag, srcStackName, srcHappyClient.HappyConfig.GetEnv())
	}

	// make sure the target builder is using the tags that were just pulled
	targetHappyClient.ArtifactBuilder.WithTags([]string{srcTag})
	return nil
}

func validate(validations ...validation) error {
	for _, validation := range validations {
		logrus.Debugf("Running validation: %s", runtime.FuncForPC(reflect.ValueOf(validation).Pointer()).Name())
		err := validation()
		if err != nil {
			return errors.Wrap(err, "unable to validate the environment")
		}
	}
	return nil
}

func makeWaitOptions(stackName string, happyConfig *config.HappyConfig, backend *backend.Backend) waitoptions.WaitOptions {
	taskOrchestrator := orchestrator.NewOrchestrator().WithHappyConfig(happyConfig).WithBackend(backend)
	return waitoptions.WaitOptions{
		StackName:    stackName,
		Orchestrator: taskOrchestrator,
		Services:     happyConfig.GetServices(),
	}
}

func verifyTFEBacklog(ctx context.Context, workspaceRepo workspace_repo.WorkspaceRepoIface) error {
	if !diagnostics.IsInteractiveContext(ctx) {
		// When you're not interactive, no point in measuring the backlog size
		return nil
	}
	backlogSize, _, err := workspaceRepo.EstimateBacklogSize(ctx)
	if err != nil {
		return errors.Wrap(err, "error estimating TFE backlog")
	}
	if backlogSize < 2 {
		logrus.Debug("There is no TFE backlog, proceeding.")
	} else if backlogSize < 20 {
		logrus.Debugf("TFE backlog is only %d runs long, proceeding.", backlogSize)
	} else {
		proceed := false
		prompt := &survey.Confirm{Message: fmt.Sprintf("TFE backlog is %d runs long, it might take a while to clear out. Do you want to wait? ", backlogSize)}
		err = survey.AskOne(prompt, &proceed)
		if err != nil {
			return errors.Wrapf(err, "failed to ask for confirmation")
		}

		if !proceed {
			return err
		}
	}
	return nil
}

// TODO: Convert to validations once HappyClient is implemented in all commands
func stackExists(stacks map[string]*stackservice.Stack, stackName string) (*stackservice.Stack, bool) {
	stack, ok := stacks[stackName]
	return stack, ok
}

func serviceExists(happyConfig *config.HappyConfig, serviceName string) bool {
	for _, s := range happyConfig.GetServices() {
		if s == serviceName {
			return true
		}
	}
	return false
}
