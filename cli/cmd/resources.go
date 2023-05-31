package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
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
	rootCmd.AddCommand(resourcesCmd)
	config.ConfigureCmdWithBootstrapConfig(resourcesCmd)
}

var resourcesCmd = &cobra.Command{
	Use:          "resources",
	Short:        "Get stack resources",
	Long:         "Get stack resources in environment '{env}'",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1)),
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

		tfeUrl := b.Conf().GetTfeUrl()
		tfeOrg := b.Conf().GetTfeOrg()

		workspaceRepo := workspace_repo.NewWorkspaceRepo(tfeUrl, tfeOrg)
		stackSvc := stackservice.NewStackService().WithHappyConfig(happyConfig).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		stack, ok := stacks[stackName]
		if !ok {
			return errors.Errorf("stack '%s' not found in environment '%s'", stackName, happyConfig.GetEnv())
		}

		logrus.Infof("Retrieving stack '%s' from environment '%s'", stackName, happyConfig.GetEnv())

		taskOrchestrator := orchestrator.NewOrchestrator().WithHappyConfig(happyConfig).WithBackend(b)
		resources, err := taskOrchestrator.GetResources(ctx, stack)
		if err != nil {
			return errors.Wrapf(err, "error retrieving resources for stack '%s'", stackName)
		}

		cb, err := b.GetComputeBackend(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to connect to a compute backend")
		}

		computeResources, err := cb.GetResources(ctx, stackName)
		if err != nil {
			return errors.Wrapf(err, "error retrieving compute level resources for stack '%s'", stackName)
		}
		resources = append(resources, computeResources...)

		printer := output.NewPrinter(OutputFormat)

		err = printer.PrintResources(ctx, resources)
		if err != nil {
			return err
		}

		return nil
	},
}
