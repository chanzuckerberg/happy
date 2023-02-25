package cmd

import (
	"context"

	"github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/cli/pkg/workspace_repo"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func initializeHappyClients(cmd *cobra.Command, sliceName, tag string, createTag, dryRun bool) (
	*config.HappyConfig,
	*stackservice.StackService,
	artifact_builder.ArtifactBuilderIface,
	*backend.Backend,
	error) {
	bootstrapConfig, err := config.NewBootstrapConfig(cmd)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	ctx := cmd.Context()
	awsBackend, err := backend.NewAWSBackend(ctx, happyConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	builderConfig := artifact_builder.
		NewBuilderConfig().
		WithBootstrap(bootstrapConfig).
		WithHappyConfig(happyConfig).
		WithDryRun(dryRun)
	ab, err := configureArtifactBuilder(ctx, sliceName, tag, createTag, dryRun, builderConfig, happyConfig, awsBackend)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	workspaceRepo := createWorkspaceRepo(dryRun, awsBackend)
	stackService := stackservice.NewStackService().
		WithBackend(awsBackend).
		WithWorkspaceRepo(workspaceRepo)

	return happyConfig, stackService, ab, awsBackend, nil
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
	builderConfig *artifact_builder.BuilderConfig,
	happyConfig *config.HappyConfig,
	backend *backend.Backend) (artifact_builder.ArtifactBuilderIface, error) {
	var err error
	if sliceName != "" {
		slice, err := happyConfig.GetSlice(sliceName)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to find the slice %s", sliceName)
		}
		builderConfig.WithProfile(slice.Profile)
	}

	// if creating tag and none specified, generate the default tag
	if createTag && (tag == "") {
		tag, err = backend.GenerateTag(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to generate tag")
		}
	}

	return artifact_builder.NewArtifactBuilder(dryRun).
			WithConfig(builderConfig).
			WithBackend(backend).
			WithTags([]string{tag}),
		nil
}
