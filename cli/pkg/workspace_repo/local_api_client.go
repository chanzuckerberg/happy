package workspace_repo

import (
	"context"

	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

type LocalWorkspaceRepo struct {
	dryRun bool
}

func NewLocalWorkspaceRepo() *LocalWorkspaceRepo {
	return &LocalWorkspaceRepo{}
}

func (c *LocalWorkspaceRepo) GetWorkspace(ctx context.Context, workspaceName string) (Workspace, error) {
	config := &tfe.Config{
		Address:    "",
		Token:      "abcd1234",
		HTTPClient: nil,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create local http client")
	}

	return &TFEWorkspace{
		tfc: client,
		workspace: &tfe.Workspace{
			ID: workspaceName,
		},
		outputs: map[string]string{},
		vars:    map[string]map[string]*tfe.Variable{},
	}, nil
}

func (c *LocalWorkspaceRepo) EstimateBacklogSize(ctx context.Context) (int, map[string]int, error) {
	return 0, map[string]int{}, nil
}

func (c *LocalWorkspaceRepo) WithDryRun(dryRun bool) *LocalWorkspaceRepo {
	c.dryRun = dryRun
	return c
}
