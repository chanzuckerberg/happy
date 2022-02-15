package workspace_repo

// Add mocks as necessary
//go:generate mockgen -destination=mock_workspace.go -package=workspace_repo github.com/chanzuckerberg/happy/pkg/workspace_repo Workspace
//go:generate mockgen -destination=mock_workspace_repo.go -package=workspace_repo github.com/chanzuckerberg/happy/pkg/workspace_repo WorkspaceRepoIface
