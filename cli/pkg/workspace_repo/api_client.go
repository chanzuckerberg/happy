package workspace_repo

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/cli/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-tfe"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	tokenUnknown = iota
	tokenPresent
	tokenMissing
	tokenRefreshNeeded
)

type WorkspaceRepo struct {
	url      string
	org      string
	hostAddr string
	tfc      *tfe.Client
	dryRun   util.DryRunType
}

func NewWorkspaceRepo(url string, org string) *WorkspaceRepo {
	return &WorkspaceRepo{
		url: url,
		org: org,
	}
}

// For testing purposes only
func (c *WorkspaceRepo) WithTFEClient(tfc *tfe.Client) *WorkspaceRepo {
	c.tfc = tfc
	return c
}

func (c *WorkspaceRepo) WithDryRun(dryRun util.DryRunType) *WorkspaceRepo {
	c.dryRun = dryRun
	return c
}

func (c *WorkspaceRepo) tfeLogin() error {
	composeArgs := []string{"terraform", "login", c.hostAddr}

	tf, err := util.NewDefaultExecutor().LookPath("terraform")
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

func (c *WorkspaceRepo) getToken(hostname string) (string, error) {
	// get token from env var
	token, err := GetTfeToken(hostname)
	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve a TFE token")
	}

	return token, nil
}

func (c *WorkspaceRepo) getTfc(ctx context.Context) (*tfe.Client, error) {
	if c.tfc == nil {
		defer diagnostics.AddTfeRunInfoUrl(ctx, c.url)
		u, err := url.Parse(c.url)
		if err != nil {
			return nil, err
		}
		hostAddr := u.Hostname()
		c.hostAddr = hostAddr

		tfc, err := c.enforceClient(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create a TFE client")
		}
		c.tfc = tfc
	}

	return c.tfc, nil
}

func (c *WorkspaceRepo) enforceClient(ctx context.Context) (*tfe.Client, error) {
	var tfc *tfe.Client
	var err error
	var errs *multierror.Error
	state := tokenUnknown

	var token string
	tokenAttemptCounter := 0

	for tokenAttemptCounter < 3 {
		switch state {
		case tokenUnknown:
			token, err = c.getToken(c.hostAddr)
			if err != nil {
				errs = multierror.Append(errs, err)
				state = tokenMissing
				break
			}
			state = tokenPresent
		case tokenMissing:
			tokenAttemptCounter++
			err = c.tfeLogin()
			if err != nil {
				errs = multierror.Append(errs, err)
				break
			}
			state = tokenUnknown
		case tokenRefreshNeeded:
			if !diagnostics.IsInteractiveContext(ctx) {
				return nil, errors.Wrap(errs.ErrorOrNil(), "cannot refresh a TFE token in a non-interactive mode")
			}
			tokenAttemptCounter++
			logrus.Infof("Opening Browser window to %s to refresh TFE Token.", c.url)
			err = browser.OpenURL(c.url)
			if err != nil { // irrecoverable
				return nil, multierror.Append(errs, err).ErrorOrNil()
			}
			loggedIn := false
			prompt := &survey.Confirm{Message: "Did you complete the TFE login in the browser window?"}
			err = survey.AskOne(prompt, &loggedIn)
			if err != nil { // irrecoverable
				return nil, multierror.Append(errs, err).ErrorOrNil()
			}
			// at this point, let's check our token is ok
			state = tokenPresent

		case tokenPresent:
			tfeConfig := &tfe.Config{
				Address: c.url,
				Token:   token,
			}
			tfc, err = tfe.NewClient(tfeConfig)
			if err != nil {
				return nil, errors.Wrap(err, "error creating the TFE client")
			}
			_, err = tfc.Organizations.List(ctx, &tfe.OrganizationListOptions{})
			if err != nil {
				errs = multierror.Append(errs, err)
				state = tokenRefreshNeeded
				break
			}
			return tfc, nil
		}
	}

	return nil, errors.Wrap(errs.ErrorOrNil(), "exhausted the max number of attempts to create a TFE client in interactive mode")
}

func (c *WorkspaceRepo) Stacks() ([]string, error) {
	return []string{}, nil
}

func (c *WorkspaceRepo) GetWorkspace(ctx context.Context, workspaceName string) (Workspace, error) {
	client, err := c.getTfc(ctx)
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

func (c *WorkspaceRepo) EstimateBacklogSize(ctx context.Context) (int, map[string]int, error) {
	backlog := map[string]int{}
	client, err := c.getTfc(ctx)
	if err != nil {
		return 0, backlog, err
	}
	page := 0
	count := 0

	for {
		options := tfe.AdminRunsListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: page,
				PageSize:   100,
			},
			RunStatus: strings.Join([]string{string(tfe.RunApplying), string(tfe.RunConfirmed), string(tfe.RunCostEstimating), string(tfe.RunPlanning), string(tfe.RunPolicyChecking)}, ","),
			Include:   []tfe.AdminRunIncludeOpt{tfe.AdminRunWorkspace, tfe.AdminRunWorkspaceOrg, tfe.AdminRunWorkspaceOrgOwners},
		}
		adminRuns, err := client.Admin.Runs.List(ctx, &options)
		if err != nil {
			if errors.Is(err, tfe.ErrResourceNotFound) {
				// User does not have access to admin API
				return 0, backlog, nil
			}
			return 0, backlog, errors.Wrapf(err, "Unable to estimate the size of TFE backlog")
		}
		for _, run := range adminRuns.Items {
			if run.Workspace != nil && run.Workspace.Organization != nil {
				key := fmt.Sprintf("%s (%s)", run.Workspace.Organization.Name, run.Status)
				backlog[key] = backlog[key] + 1
			}
		}

		count += len(adminRuns.Items)
		if adminRuns.Pagination.NextPage == 0 || page > 10 {
			break
		}

		page = adminRuns.NextPage
	}

	return count, backlog, nil
}
