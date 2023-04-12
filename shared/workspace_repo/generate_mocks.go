package workspace_repo

// Add mocks as necessary
//go:generate mockgen -destination=mock_workspace_test.go -package=workspace_repo github.com/chanzuckerberg/happy/shared/workspace_repo Workspace
//go:generate mockgen -destination=mock_workspace_repo_test.go -package=workspace_repo github.com/chanzuckerberg/happy/shared/workspace_repo WorkspaceRepoIface
