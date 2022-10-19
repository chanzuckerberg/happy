package workspace_repo

import (
	"context"

	"github.com/chanzuckerberg/happy/pkg/cli/options"
	"github.com/chanzuckerberg/happy/pkg/cli/util"
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
	GetLatestConfigVersionID() (string, error)
	Run(isDestroy bool, dryRun util.DryRunType) error
	SetVars(key string, value string, description string, sensitive bool) error
	RunConfigVersion(configVersionId string, isDestroy bool, dryRun util.DryRunType) error
	Wait(ctx context.Context, dryRun util.DryRunType) error
	WaitWithOptions(ctx context.Context, waitOptions options.WaitOptions, dryRun util.DryRunType) error
	ResetCache()
	GetTags() (map[string]string, error)
	GetWorkspaceId() string
	GetOutputs() (map[string]string, error)
	GetCurrentRunStatus() string
	UploadVersion(targzFilePath string, dryRun util.DryRunType) (string, error)
	SetOutputs(map[string]string)          // For testing purposes only
	SetClient(tfc *tfe.Client)             // For testing purposes only
	SetWorkspace(workspace *tfe.Workspace) // For testing purposes only
	HasState(ctx context.Context) (bool, error)
	DiscardRun(ctx context.Context, runID string) error
}
