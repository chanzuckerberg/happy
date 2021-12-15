package mocks

// Add mocks as necessary
//go:generate mockgen -destination=mock_backend.go -package=mocks github.com/chanzuckerberg/happy-deploy/pkg/backend ParamStoreBackend
//go:generate mockgen -destination=mock_workspace.go -package=mocks github.com/chanzuckerberg/happy-deploy/pkg/workspace_repo Workspace
//go:generate mockgen -destination=mock_workspace_repo.go -package=mocks github.com/chanzuckerberg/happy-deploy/pkg/workspace_repo WorkspaceRepoIface
//go:generate mockgen -destination=mock_dir_processor.go -package=mocks github.com/chanzuckerberg/happy-deploy/pkg/util DirProcessor
