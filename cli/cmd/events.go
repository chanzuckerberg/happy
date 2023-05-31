package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(eventsCmd)
	config.ConfigureCmdWithBootstrapConfig(eventsCmd)
}

var eventsCmd = &cobra.Command{
	Use:          "events",
	Short:        "Show stack events",
	Long:         "Showing stack events in environment '{env}'",
	PreRunE:      cmd.Validate(cobra.ExactArgs(1)),
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
			checklist.AwsInstalled,
		)
	},
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
		stackSvc := stackservice.NewStackService().WithHappyConfig(happyConfig).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		_, ok := stacks[stackName]
		if !ok {
			return errors.Errorf("stack '%s' not found in environment '%s'", stackName, happyConfig.GetEnv())
		}

		logrus.Infof("Retrieving stack '%s' events from environment '%s'", stackName, happyConfig.GetEnv())

		err = b.GetEvents(ctx, stackName, happyConfig.GetServices())
		if err != nil {
			return err
		}

		return nil
	},
}
