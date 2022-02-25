package workspace_repo

import (
	"context"
	"net/url"
	"os"
	"os/exec"

	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

const (
	tokenUnknown = 0
	tokenPresent = 1
	tokenMissing = 2
)

type WorkspaceRepo struct {
	url      string
	org      string
	hostAddr string
	ctx      context.Context
	tfc      *tfe.Client
}

func NewWorkspaceRepo(ctx context.Context, url string, org string) (*WorkspaceRepo, error) {
	// TODO do a check if see if token for the workspace repo (TFE) has expired
	return &WorkspaceRepo{
		url: url,
		org: org,
	}, nil
}

func (c *WorkspaceRepo) tfeLogin() error {
	composeArgs := []string{"terraform", "login", c.hostAddr}

	tf, err := exec.LookPath("terraform")
	if err != nil {
		return errors.Wrap(err, "terraform not in path")
	}

	cmd := &exec.Cmd{
		Path:   tf,
		Args:   composeArgs,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	err = util.NewDefaultExecutor().Run(cmd)
		return errors.Wrap(err, "unable to execute terraform")
}

func (c *WorkspaceRepo) getToken(ctx context.Context, hostname string) (string, error) {
	// get token from env var
	token, err := GetTfeToken(ctx, hostname)
	if err != nil {
		return "", errors.Wrap(err, "please set env var TFE_TOKEN")
	}

	return token, nil
}

func (c *WorkspaceRepo) getTfc(ctx context.Context) (*tfe.Client, error) {
	if c.tfc == nil {
		u, err := url.Parse(c.url)
		if err != nil {
			return nil, err
		}
		hostAddr := u.Hostname()
		c.hostAddr = hostAddr

		tfc, err := c.enforceClient(ctx)
		if err != nil {
			c.tfc = tfc
		} else {
			return nil, errors.Wrap(err, "unable to create a TFE client")
		}
	}

	return c.tfc, nil
}

func (c *WorkspaceRepo) enforceClient(ctx context.Context) (*tfe.Client, error) {
	var tfc *tfe.Client
	var err error
	var errs *multierror.Error
	state := tokenUnknown

	var token string
	tokenPresentCounter := 0

	for tokenPresentCounter < 3 {
		switch state {
		case tokenUnknown:
			token, err = c.getToken(ctx, c.hostAddr)
			if err != nil {
			        errs = multierror.Append(errs, err)
			        state = tokenMissing
                                break
			} 
			state = tokenPresent
		case tokenMissing:
			err = c.tfeLogin()
			if err =! nil {
				errs = multierror.Append(errs, err)
				break
			}
			state = tokenUnknown
			
		case tokenPresent:
			tfeConfig := &tfe.Config{
				Address: c.url,
				Token:   token,
			}
			tfc, err = tfe.NewClient(tfeConfig)
			if err != nil {
				return nil, errors.Wrap(err, "error creating the TFE client")
			}
			_, err = tfc.Organizations.List(ctx, tfe.OrganizationListOptions{})

			if err == nil {
				return tfc, nil
			} else {
				state = tokenMissing
			}
		}
	}
	return nil, errors.Wrap(err, "exhausted the max number of attempts to create a TFE client")
}

func (c *WorkspaceRepo) Stacks() ([]string, error) {
	return []string{}, nil
}

func (c *WorkspaceRepo) GetWorkspace(workspaceName string) (Workspace, error) {
	client, err := c.getTfc(c.ctx)
	if err != nil {
		return nil, err
	}

	ws, err := client.Workspaces.Read(context.Background(), c.org, workspaceName)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read workspace %s", workspaceName)
	}

	tfeWorkspace := &TFEWorkspace{
		tfc:       client,
		workspace: ws,
	}
	// Make sure we populate all variables in the workspace
	_, err = tfeWorkspace.getVars()
	return tfeWorkspace, err
}
