package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	ab "github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/cli/pkg/diagnostics"
	waitoptions "github.com/chanzuckerberg/happy/cli/pkg/options"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/cli/pkg/workspace_repo"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type HappyClient struct {
	HappyConfig     *config.HappyConfig
	StackService    *stackservice.StackService
	ArtifactBuilder ab.ArtifactBuilderIface
	Tag             string
	StackTags       map[string]string
	AWSBackend      *backend.Backend
}

func makeHappyClient(cmd *cobra.Command, sliceName, stackName, tag string, createTag, dryRun bool) (*HappyClient, error) {
	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
	if err != nil {
		return nil, err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, err
	}
	ctx := cmd.Context()
	awsBackend, err := backend.NewAWSBackend(ctx, happyConfig)
	if err != nil {
		return nil, err
	}
	builderConfig := ab.
		NewBuilderConfig().
		WithBootstrap(bootstrapConfig).
		WithHappyConfig(happyConfig)

	builderConfig.DryRun = dryRun
	builderConfig.StackName = stackName
	ab, tag, stackTags, err := configureArtifactBuilder(ctx, sliceName, tag, createTag, dryRun, builderConfig, happyConfig, awsBackend)
	if err != nil {
		return nil, err
	}
	workspaceRepo := createWorkspaceRepo(dryRun, awsBackend)
	stackService := stackservice.NewStackService().
		WithBackend(awsBackend).
		WithWorkspaceRepo(workspaceRepo)

	return &HappyClient{
		HappyConfig:     happyConfig,
		StackService:    stackService,
		ArtifactBuilder: ab,
		Tag:             tag,
		StackTags:       stackTags,
		AWSBackend:      awsBackend,
	}, nil
}

func createWorkspaceRepo(isDryRun bool, backend *backend.Backend) workspace_repo.WorkspaceRepoIface {
	if util.IsLocalstackMode() {
		return workspace_repo.NewLocalWorkspaceRepo().WithDryRun(isDryRun)
	}
	url := backend.Conf().GetTfeUrl()
	org := backend.Conf().GetTfeOrg()
	return workspace_repo.NewWorkspaceRepo(url, org).WithDryRun(isDryRun)
}

func configureArtifactBuilder(
	ctx context.Context,
	sliceName, tag string,
	createTag, dryRun bool,
	builderConfig *ab.BuilderConfig,
	happyConfig *config.HappyConfig,
	backend *backend.Backend) (ab.ArtifactBuilderIface, string, map[string]string, error) {
	var err error
	if sliceName != "" {
		slice, err := happyConfig.GetSlice(sliceName)
		if err != nil {
			return nil, "", nil, errors.Wrapf(err, "unable to find the slice %s", sliceName)
		}
		builderConfig.Profile = slice.Profile
	}

	// if creating tag and none specified, generate the default tag
	if createTag && (tag == "") {
		tag, err = backend.GenerateTag(ctx)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "unable to generate tag")
		}
	}

	stackTags := map[string]string{}
	if sliceName != "" {
		serviceImages, err := builderConfig.GetBuildServicesImage(ctx)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "unable to get build service images")
		}

		for service := range serviceImages {
			stackTags[service] = tag
		}
	}

	return ab.NewArtifactBuilder(dryRun).
		WithConfig(builderConfig).
		WithBackend(backend).
		WithTags([]string{tag}), tag, stackTags, nil
}

type validation func() error

func validateImageExists(ctx context.Context, createTag, skipCheckTag bool, ab ab.ArtifactBuilderIface) validation {
	return func() error {
		if skipCheckTag {
			return nil
		}

		if createTag {
			// if we build and push and it succeeds, we know that the image exists
			return ab.BuildAndPush(ctx)
		}

		if len(ab.GetTags()) == 0 {
			return errors.Errorf("no tags have been assigned")
		}

		exists, err := ab.CheckImageExists(ctx, ab.GetTags()[0])
		if err != nil {
			return errors.Wrapf(err, "error checking if tag %s existed", ab.GetTags()[0])
		}
		if !exists {
			return errors.Errorf("image tag does not exist: '%s'", ab.GetTags()[0])
		}

		return nil
	}
}
func validateTFEBackLog(ctx context.Context, isDryRun bool, awsBackend *backend.Backend) validation {
	return func() error {
		return verifyTFEBacklog(ctx, createWorkspaceRepo(isDryRun, awsBackend))
	}
}

func validateGitTree(projectRoot string) validation {
	return func() error {
		return util.ValidateGitTree(projectRoot)
	}
}

func validateStackNameAvailable(ctx context.Context, stackService *stackservice.StackService, stackName string, force bool) validation {
	return func() error {
		if force {
			return nil
		}

		_, err := stackService.GetStack(ctx, stackName)
		if err != nil {
			return nil
		}

		return errors.Wrap(err, "the stack name is already taken")
	}
}

func validate(validations ...validation) error {
	for _, validation := range validations {
		err := validation()
		if err != nil {
			return errors.Wrap(err, "unable to validate the environment")
		}
	}
	return nil
}

func makeWaitOptions(stackName string, backend *backend.Backend) waitoptions.WaitOptions {
	taskOrchestrator := orchestrator.NewOrchestrator().WithBackend(backend)
	return waitoptions.WaitOptions{
		StackName:    stackName,
		Orchestrator: taskOrchestrator,
		Services:     backend.Conf().GetServices(),
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
