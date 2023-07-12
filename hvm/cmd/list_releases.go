/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/hvm/github_connector"
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
	fmt.Println("list-releases called")

	client := github_connector.NewConnectorClient()
	releases, err := client.GetHappyReleases()

	if err != nil {
		fmt.Println("An error occurred getting the release list: ", err)
		return err
	}

	for _, release := range releases {
		fmt.Println(release)
	}

	return nil
}
