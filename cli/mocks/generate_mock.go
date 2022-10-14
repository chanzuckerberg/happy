package mocks

// Add mocks as necessary
//go:generate mockgen -destination=mock_workspace.go -package=mocks github.com/chanzuckerberg/happy/pkg/workspace_repo Workspace
//go:generate mockgen -destination=mock_workspace_repo.go -package=mocks github.com/chanzuckerberg/happy/pkg/workspace_repo WorkspaceRepoIface
//go:generate mockgen -destination=mock_dir_processor.go -package=mocks github.com/chanzuckerberg/happy/pkg/util DirProcessor
//go:generate mockgen -destination=mock_stack.go -package=mocks github.com/chanzuckerberg/happy/pkg/stack_mgr StackIface
//go:generate mockgen -destination=mock_stack_service.go -package=mocks github.com/chanzuckerberg/happy/pkg/stack_mgr StackServiceIface
//go:generate mockgen -destination=mock_artifact_builder.go -package=mocks github.com/chanzuckerberg/happy/pkg/artifact_builder ArtifactBuilderIface
