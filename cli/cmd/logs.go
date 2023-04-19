package cmd

import (
	"os"
	"time"

	"github.com/chanzuckerberg/happy/cli/pkg/cmd"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	since string
	//follow     bool
	outputFile string
)

func init() {
	rootCmd.AddCommand(logsCmd)
	config.ConfigureCmdWithBootstrapConfig(logsCmd)

	logsCmd.Flags().StringVar(&since, "since", "1h", "Length of time to look back in logs, ex. 10s, 5m, 24h.")
	//logsCmd.Flags().BoolVar(&follow, "follow", false, "Specify if the logs should be streamed")
	logsCmd.Flags().StringVar(&outputFile, "output", "", "Specify if the logs should be output to a file")
}

var logsCmd = &cobra.Command{
	Use:          "logs STACK_NAME SERVICE",
	Short:        "Print logs",
	Long:         "Print the logs of a service (frontend, backend, upload, migrations)",
	SilenceUsage: true,
	RunE:         runLogs,
	PreRunE:      cmd.Validate(cobra.ExactArgs(2), cmd.IsStackNameDNSCharset),
}

func runLogs(cmd *cobra.Command, args []string) error {
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

	url := b.Conf().GetTfeUrl()
	org := b.Conf().GetTfeOrg()

	workspaceRepo := workspace_repo.NewWorkspaceRepo(url, org)
	stackSvc := stackservice.NewStackService().WithHappyConfig(happyConfig).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

	stacks, err := stackSvc.GetStacks(ctx)
	if err != nil {
		return err
	}
	stackExists := func() bool {
		for _, stack := range stacks {
			if stack.Name == stackName {
				return true
			}
		}
		return false
	}()
	if !stackExists {
		return errors.Errorf("stack %s doesn't exist for env %s", stackName, happyConfig.GetEnv())
	}
	serviceExists := func() bool {
		for _, s := range happyConfig.GetServices() {
			if s == serviceName {
				return true
			}
		}
		return false
	}()
	if !serviceExists {
		return errors.Errorf("service %s doesn't exist for env %s. available services: %+v", serviceName, happyConfig.GetEnv(), happyConfig.GetServices())
	}

	opts := []util.PrintOption{}
	if outputFile != "" {
		writer, err := os.Create(outputFile)
		if err != nil {
			return errors.Wrap(err, "error opening file for logging")
		}
		defer writer.Close()
		opts = append(opts, util.WithWriter(writer))
	}

	if since != "" {
		duration, err := time.ParseDuration(since)
		if err != nil {
			return errors.Wrapf(err, "unable to parse the 'since' param %s", since)
		}
		opts = append(opts, util.WithSince(util.GetStartTime(ctx).Add(-duration).UnixMilli()))
	}

	return b.PrintLogs(ctx, stackName, serviceName, opts...)
}
