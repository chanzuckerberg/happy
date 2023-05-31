package cmd

import (
	"io"

	"github.com/chanzuckerberg/happy/cli/pkg/output"
	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type StructuredListResult struct {
	Error  string
	Stacks []stackservice.StackInfo
}

var listAll bool

func init() {
	RootCmd.AddCommand(listCmd)
	config.ConfigureCmdWithBootstrapConfig(listCmd)
	listCmd.Flags().StringVar(&OutputFormat, "output", "text", "Output format. One of: json, yaml, or text. Defaults to text, which is the only interactive mode.")
	listCmd.Flags().BoolVar(&listAll, "all", false, "List all stacks, not just those belonging to this app")
}

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List stacks",
	Long:         "Listing stacks in environment '{env}'",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if OutputFormat != "text" {
			logrus.SetOutput(io.Discard)
		}
		happyClient, err := makeHappyClient(cmd, sliceName, "", []string{}, false)
		if err != nil {
			return errors.Wrap(err, "unable to initialize the happy client")
		}

		stackInfos, err := happyClient.StackService.CollectStackInfo(cmd.Context(), listAll, happyClient.HappyConfig.App())
		if err != nil {
			return errors.Wrap(err, "unable to collect stack info")
		}

		printer := output.NewPrinter(OutputFormat)
		err = printer.PrintStacks(cmd.Context(), stackInfos)
		if err != nil {
			return errors.Wrap(err, "unable to print stacks")
		}

		return nil
	},
}
