package mocks

// Add mocks as necessary
//go:generate mockgen -destination=mock_workspace.go -package=mocks github.com/chanzuckerberg/happy/shared/workspace_repo Workspace
//go:generate mockgen -destination=mock_workspace_repo.go -package=mocks github.com/chanzuckerberg/happy/shared/workspace_repo WorkspaceRepoIface
//go:generate mockgen -destination=mock_artifact_builder.go -package=mocks github.com/chanzuckerberg/happy/cli/pkg/artifact_builder ArtifactBuilderIface
