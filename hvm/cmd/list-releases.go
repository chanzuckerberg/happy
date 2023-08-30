package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/shared/githubconnector"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

)

// useCmd represents the use command
var listRelasesCommand = &cobra.Command{
	Use: "list-releases [org] [project]",

	Short: "Get list of available releases",
	Long:  `List latest releases for a project. May not be comprehensive.`,
	RunE:  listReleases,
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
		return errors.Wrap(err, "getting release list")
	}

	for _, release := range releases {
		fmt.Println(release.Version)
	}

	return nil
}
