package util

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const GithubGraphQLEndpoint = "https://api.github.com/graphql"

type githubDeploymentsResponse struct {
	Repository struct {
		Deployments struct {
			Nodes []struct {
				CommitOid string `json:"commitOid"`
				Statuses  struct {
					Nodes []struct {
						State     string    `json:"state"`
						UpdatedAt time.Time `json:"updatedAt"`
					} `json:"nodes"`
				} `json:"statuses"`
			} `json:"nodes"`
		} `json:"deployments"`
	} `json:"repository"`
}

func GetLatestSuccessfulDeployment(ctx context.Context, endpoint string, token string, stage string, owner string, repo string) (string, error) {
	query := `query ($repo_owner: String!, $repo_name: String!, $deployment_env: String!) {
		repository(owner: $repo_owner, name: $repo_name) {
		  deployments(environments: [$deployment_env], last: 50) {
			nodes {
			  commitOid
			  statuses(first: 100) {
				nodes {
				  state
				  updatedAt
				}
			  }
			}
		  }
		}
	  }`
	req := graphql.NewRequest(query)
	req.Var("repo_owner", owner)
	req.Var("repo_name", repo)
	req.Var("deployment_env", stage)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	var resp githubDeploymentsResponse

	client := graphql.NewClient(endpoint)
	if err := client.Run(ctx, req, &resp); err != nil {
		return "", errors.Wrap(err, "failed to execute a graphql request")
	}

	sha := ""
	timestamp := time.Unix(0, 0)
	for _, node := range resp.Repository.Deployments.Nodes {
		for _, status := range node.Statuses.Nodes {
			if status.State == "SUCCESS" {
				if status.UpdatedAt.After(timestamp) {
					timestamp = status.UpdatedAt
					sha = node.CommitOid
				}
			}
		}
	}
	if len(sha) > 8 {
		sha = sha[:8]
	}
	return sha, nil
}

func IsCleanGitTree(dir string) (bool, *git.Status, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logrus.Debugf("Folder %s is not under source control, or git is not installed", dir)
		return true, nil, nil
	}
	r, err := git.PlainOpen(strings.Trim(out.String(), "\n"))
	if err != nil {
		return false, nil, errors.Wrap(err, "unable to open git repository")
	}
	w, err := r.Worktree()
	if err != nil {
		return false, nil, errors.Wrap(err, "cannot get the working tree of git repository")
	}
	status, err := w.Status()
	if err != nil {
		return false, nil, errors.Wrap(err, "cannot get the status of the working tree")
	}
	return status.IsClean(), &status, nil
}

func ValidateGitTree(dir string) error {
	isClean, status, err := IsCleanGitTree(dir)
	if err != nil {
		return err
	}
	if !isClean {
		var dirtyFiles string
		for k := range *status {
			dirtyFiles += fmt.Sprintf("\t- %s\n", k)
		}
		logrus.Warnf("IN THE FUTURE, THIS WARNING WILL PREVENT UPDATES/CREATIONS TO STACKS\ngit tree is dirty; please commit or discard all changes below:\n%s", dirtyFiles)
	}

	return nil
}
