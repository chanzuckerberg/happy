package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/output"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
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
	Short:        "Get stack",
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

		b, err := backend.NewAWSBackend(ctx, happyConfig.GetEnvironmentContext())
		if err != nil {
			return err
		}

		tfeUrl := b.Conf().GetTfeUrl()
		tfeOrg := b.Conf().GetTfeOrg()

		workspaceRepo := workspace_repo.NewWorkspaceRepo(tfeUrl, tfeOrg)
		stackSvc := stackservice.NewStackService(happyConfig).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		stack, ok := stacks[stackName]
		if !ok {
			return errors.Errorf("stack '%s' not found in environment '%s'", stackName, happyConfig.GetEnv())
		}

		logrus.Infof("Retrieving stack '%s' from environment '%s'", stackName, happyConfig.GetEnv())

		tablePrinter := util.NewTablePrinter()

		stackInfo, err := stack.GetStackInfo(ctx, stackName)
		if err != nil {
			return errors.Wrapf(err, "error retrieving stack '%s'", stackName)
		}

		tablePrinter.Print(output.Stack2Console(*stackInfo))

		backlogSize, backlog, err := workspaceRepo.EstimateBacklogSize(ctx)
		if err != nil {
			return errors.Wrap(err, "error estimating TFE backlog")
		}

		tablePrinter = util.NewTablePrinter()

		tablePrinter.AddSimpleRow("Environment", bootstrapConfig.Env)
		tablePrinter.AddSimpleRow("TFE", "")
		tablePrinter.AddSimpleRow("  Environment Workspace", fmt.Sprintf("%s/app/%s/workspaces/env-%s", tfeUrl, tfeOrg, bootstrapConfig.Env))
		tablePrinter.AddSimpleRow("  Stack Workspace", fmt.Sprintf("%s/app/%s/workspaces/%s-%s", tfeUrl, tfeOrg, bootstrapConfig.Env, stackName))

		if len(backlog) > 0 {
			tablePrinter.AddSimpleRow("  Backlog size", fmt.Sprintf("%d outstanding runs", backlogSize))
			for k, v := range backlog {
				tablePrinter.AddSimpleRow("", fmt.Sprintf("%s->%d", k, v))
			}
		}

		tablePrinter.AddSimpleRow("AWS", "")
		tablePrinter.AddSimpleRow("  Account ID", fmt.Sprintf("[%s]", b.GetAWSAccountID()))
		tablePrinter.AddSimpleRow("  Region", b.GetAWSRegion())
		tablePrinter.AddSimpleRow("  Profile", b.GetAWSProfile())

		for _, serviceName := range happyConfig.GetServices() {
			tablePrinter.AddSimpleRow("Service", serviceName)

			descriptor, err := b.Describe(ctx, stackName, serviceName)
			if err != nil {
				return errors.Errorf("error describing service %s", serviceName)
			}
			tablePrinter.AddSimpleRow("  Compute", descriptor.Compute)
			for k, v := range descriptor.Params {
				tablePrinter.AddSimpleRow(fmt.Sprintf("  %s", k), v)
			}
		}
		tablePrinter.Flush()

		return nil
	},
}
