package workspace_repo

import (
	"github.com/chanzuckerberg/happy/pkg/options"
	tfe "github.com/hashicorp/go-tfe"
)

type WorkspaceRepoIface interface {
	GetWorkspace(workspaceName string) (Workspace, error)
}

type Workspace interface {
	GetWorkspaceID() string
	WorkspaceName() string
	GetCurrentRunID() string
	GetLatestConfigVersionID() (string, error)
	Run(isDestroy bool) error
	SetVars(key string, value string, description string, sensitive bool) error
	RunConfigVersion(configVersionId string, isDestroy bool) error
	Wait() error
	WaitWithOptions(waitOptions options.WaitOptions) error
	ResetCache()
	GetTags() (map[string]string, error)
	GetWorkspaceId() string
	GetOutputs() (map[string]string, error)
	GetCurrentRunStatus() string
	UploadVersion(targzFilePath string) (string, error)
	SetClient(tfc *tfe.Client)             // For testing purposes only
	SetWorkspace(workspace *tfe.Workspace) // For testing purposes only
	SetCurrentRun(currentRun *tfe.Run)     // For testing purposes only
}
