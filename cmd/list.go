package cmd

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

var outputFormat string

type StructuredListResult struct {
	Error  string
	Stacks []stackservice.StackInfo
}

func init() {
	rootCmd.AddCommand(listCmd)
	config.ConfigureCmdWithBootstrapConfig(listCmd)
	listCmd.Flags().StringVar(&outputFormat, "output", "text", "Output format. One of: json, yaml, or text. Defaults to text, which is the only interactive mode.")
}

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "list stacks",
	Long:         "Listing stacks in environment '{env}'",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputFormat != "text" {
			logrus.SetOutput(ioutil.Discard)
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

		b, err := backend.NewAWSBackend(ctx, happyConfig)
		if err != nil {
			return err
		}

		url := b.Conf().GetTfeUrl()
		org := b.Conf().GetTfeOrg()

		workspaceRepo := workspace_repo.NewWorkspaceRepo(url, org).WithInteractive(Interactive)
		stackSvc := stackservice.NewStackService().WithBackend(b).WithWorkspaceRepo(workspaceRepo)

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
				if !Interactive {
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

		if outputFormat == "text" {
			logrus.Infof("listing stacks in environment '%s'", happyConfig.GetEnv())

			headings := []string{"Name", "Owner", "Tags", "Status", "URLs", "LastUpdated"}
			tablePrinter := util.NewTablePrinter(headings)

			for _, stackInfo := range stackInfos {
				tablePrinter.AddRow(stackInfo.Name, stackInfo.Owner, stackInfo.Tag, stackInfo.Status, stackInfo.Url, stackInfo.LastUpdated)
			}
			tablePrinter.Print()
			return nil
		}

		if outputFormat == "json" {
			b, err := json.Marshal(stackInfos)
			if err != nil {
				return err
			}
			printOutput(string(b))
		}

		if outputFormat == "yaml" {
			b, err := yaml.Marshal(stackInfos)
			if err != nil {
				return err
			}
			printOutput(string(b))
		}

		return nil
	},
}
