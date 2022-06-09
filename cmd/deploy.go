package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd)
	config.ConfigureCmdWithBootstrapConfig(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:          "deploy",
	Short:        "deploy deployment_stage",
	Long:         "Get a sha of the last successful deployment to deployment_stage",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		stage := args[0]

		token, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return errors.New("please set GITHUB_TOKEN environment variable")
		}

		sha, err := getLatestSuccessfulDeployment(ctx, token, stage)
		if err != nil {
			return errors.Wrap(err, "failed to get last successful deployment")
		}

		fmt.Printf("%s\n", sha)

		return nil
	},
}

func getLatestSuccessfulDeployment(ctx context.Context, token string, stage string) (string, error) {
	client := graphql.NewClient("https://api.github.com/graphql")

	req := graphql.NewRequest(`
    query($repo_owner:String!, $repo_name:String!, $deployment_env:String!) {
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
	}
	`)

	// set any variables
	// req.Var("repo_owner", "chanzuckerberg")
	// req.Var("repo_name", "czgenepi")
	// req.Var("deployment_env", "stage")

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	// run it and capture the response
	var respData map[string]interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		return "", errors.Wrap(err, "failed to execute a graphql request")
	}

	response, err := json.Marshal(respData)
	if err != nil {
		return "", nil
	}

	log.Infof("Response: %v", string(response))

	return "", nil
}
