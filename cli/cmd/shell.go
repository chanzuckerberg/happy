package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/orchestrator"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	stackservice "github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var containerName string
var shellCommand string

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVarP(&containerName, "container", "c", "", "Container name")
	shellCmd.Flags().StringVar(&shellCommand, "command", "", "Command to run in the container")
	config.ConfigureCmdWithBootstrapConfig(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:          "shell STACK_NAME SERVICE",
	Aliases:      []string{"exec", "sh", "bash"},
	Short:        "Execute into a container",
	Long:         "Execute into a running service task container",
	SilenceUsage: true,
	PreRunE: cmd.Validate(
		cobra.ExactArgs(2),
		cmd.IsStackNameDNSCharset,
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()
			return util.ValidateEnvironment(cmd.Context(),
				checklist.AwsInstalled,
				checklist.AwsSessionManagerPluginInstalled,
			)
		},
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		stackName := args[0]
		serviceName := args[1]

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

		workspaceRepo := workspace_repo.NewWorkspaceRepo(b.Conf().GetTfeUrl(), b.Conf().GetTfeOrg())
		stackSvc := stackservice.NewStackService(happyConfig.GetEnv(), happyConfig.App()).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		_, stackExists := stackExists(stacks, stackName)
		if !stackExists {
			return errors.Errorf("stack %s doesn't exist for env %s", stackName, happyConfig.GetEnv())
		}
		serviceExists := serviceExists(happyConfig, serviceName)
		if !serviceExists {
			return errors.Errorf("service %s doesn't exist for env %s. available services: %+v", serviceName, happyConfig.GetEnv(), happyConfig.GetServices())
		}

		return orchestrator.NewOrchestrator().WithHappyConfig(happyConfig).WithBackend(b).Shell(ctx, stackName, serviceName, containerName, shellCommand)
	},
}
