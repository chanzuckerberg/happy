package cmd

import (
	"io"
	"sort"

	"github.com/chanzuckerberg/happy/cli/pkg/output"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type StructuredListResult struct {
	Error  string
	Stacks []stackservice.StackInfo
}

func init() {
	rootCmd.AddCommand(listCmd)
	config.ConfigureCmdWithBootstrapConfig(listCmd)
	listCmd.Flags().StringVar(&OutputFormat, "output", "text", "Output format. One of: json, yaml, or text. Defaults to text, which is the only interactive mode.")
}

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List stacks",
	Long:         "Listing stacks in environment '{env}'",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if OutputFormat != "text" {
			logrus.SetOutput(io.Discard)
		}

		ctx := cmd.Context()
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

		workspaceRepo := createWorkspaceRepo(false, b)
		stackSvc := stackservice.NewStackService(happyConfig).WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		// Iterate in order
		stackInfos := []stackservice.StackInfo{}
		stackNames := maps.Keys(stacks)
		sort.Strings(stackNames)
		for _, name := range stackNames {
			stack := stacks[name]
			stackInfo, err := stack.GetStackInfo(ctx, name)
			if err != nil {
				logrus.Warnf("Error retrieving stack %s:  %s", name, err)
				if !diagnostics.IsInteractiveContext(ctx) {
					stackInfos = append(stackInfos, stackservice.StackInfo{
						Name:    name,
						Status:  "error",
						Message: err.Error(),
					})
				}
				continue
			}
			stackInfos = append(stackInfos, *stackInfo)
		}

		logrus.Debugf("listing stacks in environment '%s'", happyConfig.GetEnv())
		printer := output.NewPrinter(OutputFormat)

		err = printer.PrintStacks(stackInfos)
		if err != nil {
			return err
		}

		return nil
	},
}
