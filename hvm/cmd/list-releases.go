package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/shared/githubconnector"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var listRelasesCommand = &cobra.Command{
	Use: "list-releases [org] [project]",

	Short: "Get list of available releases",
	Long:  ``, Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use called")
	},
	RunE: listReleases,
}

func init() {

	rootCmd.AddCommand(listRelasesCommand)
	listRelasesCommand.ArgAliases = []string{"org", "project"}
	listRelasesCommand.Args = cobra.ExactArgs(2)
}

func listReleases(cmd *cobra.Command, args []string) error {

	org := args[0]
	project := args[1]

	client := githubconnector.NewConnectorClient()
	releases, err := client.GetReleases(org, project)

	if err != nil {
		fmt.Println("An error occurred getting the release list: ", err)
		return err
	}

	for _, release := range releases {
		fmt.Println(release.Version)
	}

	return nil
}
