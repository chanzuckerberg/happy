package workspace_repo

import (
	"context"
	"net/url"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

type WorkspaceRepo struct {
	url string
	org string
	tfc *tfe.Client
}

func NewWorkspaceRepo(url string, org string) (*WorkspaceRepo, error) {
	_, err := GetTfeToken(url)
	if err != nil {
		return nil, errors.Wrap(err, "please set env var TFE_TOKEN")
	}

	// TODO do a check if see if token for the workspace repo (TFE) has expired
	return &WorkspaceRepo{
		url: url,
		org: org,
	}, nil
}

func (c *WorkspaceRepo) getToken(hostname string) (string, error) {
	// get token from env var
	token, err := GetTfeToken(hostname)
	if err != nil {
		return "", errors.Wrap(err, "please set env var TFE_TOKEN")
	}

	return token, nil
}

func (c *WorkspaceRepo) getTfc() (*tfe.Client, error) {
	if c.tfc == nil {
		// initialize the TFE client

		// TODO parse it from c.url
		// hostAddr := os.Getenv("TFE_HOST")
		u, err := url.Parse(c.url)
		if err != nil {
			return nil, err
		}
		hostAddr := u.Hostname()

		token, err := c.getToken(hostAddr)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating the service client:")
		}
		tfe_config := &tfe.Config{
			Address: c.url,
			Token:   token,
		}
		tfc, err := tfe.NewClient(tfe_config)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating the service client:")
		}

		c.tfc = tfc
	}

	return c.tfc, nil
}

func (c *WorkspaceRepo) Stacks() ([]string, error) {
	return []string{}, nil
}

func (c *WorkspaceRepo) GetWorkspace(workspaceName string) (Workspace, error) {
	client, err := c.getTfc()
	if err != nil {
		return nil, err
	}

	ws, err := client.Workspaces.Read(context.Background(), c.org, workspaceName)
	if err != nil {
		return nil, errors.Errorf("%s: %v", workspaceName, err)
	}

	return &TFEWorkspace{
		tfc:       client,
		workspace: ws,
	}, nil
}
