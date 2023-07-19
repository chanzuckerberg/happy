package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var listRelasesCommand = &cobra.Command{
	Use:   "list-releases",
	Short: "Get list of available releases",
	Long:  ``, Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use called")
	},
	RunE: listReleases,
}

func init() {
	rootCmd.AddCommand(listRelasesCommand)
}

func listReleases(cmd *cobra.Command, args []string) error {

	client := githubconnector.NewConnectorClient()
	releases, err := client.GetReleases("chanzuckerberg", "happy")

	if err != nil {
		fmt.Println("An error occurred getting the release list: ", err)
		return err
	}

	for _, release := range releases {
		fmt.Println(release.Tag)
	}

	return nil
}
