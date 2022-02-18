package cmd

import (
	"encoding/json"
	"strings"

	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	config.ConfigureCmdWithBootstrapConfig(listCmd)
}

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "list stacks",
	Long:         "Listing stacks in environment '{env}'",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		bootstrapConfig, err := config.NewBootstrapConfig()
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
		if err != nil {
			return err
		}

		b, err := backend.NewAWSBackend(ctx, happyConfig)
		if err != nil {
			return err
		}

		url := b.Conf().GetTfeUrl()
		org := b.Conf().GetTfeOrg()

		workspaceRepo, err := workspace_repo.NewWorkspaceRepo(url, org)
		if err != nil {
			return err
		}

		stackSvc := stackservice.NewStackService(b, workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		logrus.Infof("listing stacks in environment '%s'", happyConfig.GetEnv())

		// TODO look for existing package for printing table
		headings := []string{"Name", "Owner", "Tags", "Status", "URLs"}
		tablePrinter := util.NewTablePrinter(headings)

		for name, stack := range stacks {
			stackOutput, err := stack.GetOutputs()

			// TODO do not skip, just print the empty colums
			if err != nil {
				logrus.Errorf("Skipping %s due to error: %s", name, err)
				continue
			}
			url := stackOutput["frontend_url"]
			status := stack.GetStatus()
			meta, err := stack.Meta()
			if err != nil {
				return err
			}
			tag := meta.DataMap["imagetag"]
			imageTags, ok := meta.DataMap["imagetags"]
			if ok && len(imageTags) > 0 {
				var imageTagMap map[string]interface{}
				err = json.Unmarshal([]byte(imageTags), &imageTagMap)
				if err != nil {
					return err
				}
				combinedTags := []string{tag}
				for imageTag := range imageTagMap {
					combinedTags = append(combinedTags, imageTag)
				}
				tag = strings.Join(combinedTags, ", ")
			}
			tablePrinter.AddRow([]string{name, meta.DataMap["owner"], tag, status, url})
		}

		tablePrinter.Print()
		return nil
	},
}
