package cmd

import (
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
	config.ConfigureCmdWithBootstrapConfig(getCmd)
}

var getCmd = &cobra.Command{
	Use:          "get",
	Short:        "get stack",
	Long:         "Get a stack in environment '{env}'",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		stackName := args[0]

		bootstrapConfig, err := config.NewBootstrapConfig(cmd)
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(bootstrapConfig)
		if err != nil {
			return err
		}

		b, err := backend.NewAWSBackend(ctx, happyConfig)
		if err != nil {
			return err
		}

		url := b.Conf().GetTfeUrl()
		org := b.Conf().GetTfeOrg()

		workspaceRepo := workspace_repo.NewWorkspaceRepo(url, org)
		stackSvc := stackservice.NewStackService().WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		stack, ok := stacks[stackName]
		if !ok {
			return errors.Errorf("stack '%s' not found in environment '%s'", stackName, happyConfig.GetEnv())
		}

		logrus.Infof("Retrieving stack '%s' from environment '%s'", stackName, happyConfig.GetEnv())

		// TODO look for existing package for printing table
		headings := []string{"Name", "Owner", "Tags", "Status", "URLs"}
		tablePrinter := util.NewTablePrinter(headings)

		err = stack.Print(ctx, stackName, tablePrinter)

		if err != nil {
			logrus.Errorf("Error retrieving stack %s:  %s", stackName, err)
		}

		tablePrinter.Print()
		return err
	},
}
