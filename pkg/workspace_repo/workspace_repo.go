package workspace_repo

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
	ResetCache()
	GetTags() (map[string]string, error)
	GetWorkspaceId() string
	GetOutputs() (map[string]string, error)
	GetCurrentRunStatus() string
	UploadVersion(targzFilePath string) (string, error)
}
