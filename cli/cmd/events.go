package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	stackservice "github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(eventsCmd)
	config.ConfigureCmdWithBootstrapConfig(eventsCmd)
}

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Show stack events",
	Long:  "Showing stack events in environment '{env}'",
	PreRunE: cmd.Validate(
		cobra.ExactArgs(1),
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()
			return util.ValidateEnvironment(cmd.Context(),
				checklist.TerraformInstalled,
				checklist.AwsInstalled,
			)
		},
	),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		stackName := args[0]

		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		b, err := backend.NewAWSBackend(ctx, happyConfig.GetEnvironmentContext())
		if err != nil {
			return err
		}

		workspaceRepo := createWorkspaceRepo(b)
		stackSvc := stackservice.NewStackService(happyConfig.GetEnv(), happyConfig.App()).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		_, stackExists := stackExists(stacks, stackName)
		if !stackExists {
			return errors.Errorf("stack %s doesn't exist for env %s", stackName, happyConfig.GetEnv())
		}

		logrus.Infof("Retrieving stack '%s' events from environment '%s'", stackName, happyConfig.GetEnv())

		err = b.GetEvents(ctx, stackName, happyConfig.GetServices())
		if err != nil {
			return err
		}

		return nil
	},
}
