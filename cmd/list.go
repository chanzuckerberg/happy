package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	stack_service "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	config.ConfigureCmdWithBootstrapConfig(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list stacks",
	Long:  "Listing stacks in environment '{env}'",
	RunE: func(cmd *cobra.Command, args []string) error {
		bootstrapConfig, err := config.NewBootstrapConfig()
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(bootstrapConfig)
		if err != nil {
			return err
		}

		url := happyConfig.TfeUrl()
		org := happyConfig.TfeOrg()

		workspaceRepo, err := workspace_repo.NewWorkspaceRepo(url, org)
		if err != nil {
			return err
		}

		paramStoreBackend := backend.GetAwsBackend(happyConfig)
		stackSvc := stack_service.NewStackService(happyConfig, paramStoreBackend, workspaceRepo)

		stacks, err := stackSvc.GetStacks()
		if err != nil {
			return err
		}

		fmt.Printf("Listing stacks in environment '%s'\n", bootstrapConfig.GetEnv())

		// TODO look for existing package for printing table
		headings := []string{"Name", "Owner", "Tags", "Status", "URLs"}
		tablePrinter := util.NewTablePrinter(headings)

		for name, stack := range stacks {
			stackOutput, err := stack.GetOutputs()

			// TODO do not skip, just print the empty colums
			if err != nil {
				fmt.Printf("Skipping %s due to error: %s\n", name, err)
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
