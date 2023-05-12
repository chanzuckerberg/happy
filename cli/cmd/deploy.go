package cmd

import (
	"os"

	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var owner string
var repo string
var fileName string

func init() {
	rootCmd.AddCommand(deployCmd)
	config.ConfigureCmdWithBootstrapConfig(deployCmd)

	deployCmd.Flags().StringVar(&owner, "owner", "", "Repo owner (organization)")
	deployCmd.Flags().StringVar(&repo, "repo", "", "Repo name")
	deployCmd.Flags().StringVar(&fileName, "out", ".sha", "File name to output the sha to")
}

var deployCmd = &cobra.Command{
	Use:          "deploy",
	Short:        "Get a git sha of last successful deployment",
	Long:         "Get a git sha of the last successful deployment to a selected environment",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1)),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.DockerEngineRunning,
			checklist.MinDockerComposeVersion,
			checklist.DockerInstalled,
			checklist.TerraformInstalled,
			checklist.AwsInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		stage := args[0]

		token, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return errors.New("please set GITHUB_TOKEN environment variable")
		}

		sha, err := util.GetLatestSuccessfulDeployment(ctx, util.GithubGraphQLEndpoint, token, stage, owner, repo)
		if err != nil {
			return errors.Wrap(err, "failed to get last successful deployment")
		}

		log.Infof("Last successful deployment SHA: %s\n", sha)
		if len(fileName) > 0 {
			f, err := os.Create(fileName)
			if err != nil {
				return errors.Wrap(err, "cannot create a file")
			}
			defer f.Close()
			_, err = f.WriteString(sha)

			if err != nil {
				return errors.Wrap(err, "cannot write to a file")
			}
		}

		return nil
	},
}
