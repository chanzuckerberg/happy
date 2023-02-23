package workspace_repo

import (
	"context"

	"github.com/chanzuckerberg/happy/cli/pkg/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/hashicorp/go-tfe"
)

type WorkspaceRepoIface interface {
	GetWorkspace(ctx context.Context, workspaceName string) (Workspace, error)
	EstimateBacklogSize(ctx context.Context) (int, map[string]int, error)
}

type Workspace interface {
	GetWorkspaceID() string
	WorkspaceName() string
	GetCurrentRunID() string
	GetLatestConfigVersionID(ctx context.Context) (string, error)
	Run(ctx context.Context, isDestroy bool, dryRun util.DryRunType) error
	SetVars(ctx context.Context, key string, value string, description string, sensitive bool) error
	RunConfigVersion(ctx context.Context, configVersionId string, isDestroy bool, dryRun util.DryRunType) error
	Wait(ctx context.Context, dryRun util.DryRunType) error
	WaitWithOptions(ctx context.Context, waitOptions options.WaitOptions, dryRun util.DryRunType) error
	ResetCache()
	GetTags(ctx context.Context) (map[string]string, error)
	GetWorkspaceId() string
	GetOutputs(ctx context.Context) (map[string]string, error)
	GetResources(ctx context.Context) ([]util.ManagedResource, error)
	GetCurrentRunStatus(ctx context.Context) string
	UploadVersion(ctx context.Context, targzFilePath string, dryRun util.DryRunType) (string, error)
	SetOutputs(map[string]string)          // For testing purposes only
	SetClient(tfc *tfe.Client)             // For testing purposes only
	SetWorkspace(workspace *tfe.Workspace) // For testing purposes only
	HasState(ctx context.Context) (bool, error)
	DiscardRun(ctx context.Context, runID string) error
}
