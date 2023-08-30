package cmd

import (
	"fmt"
	"runtime"

	"github.com/chanzuckerberg/happy/shared/githubconnector"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [org] [project] [version]",
	Short: "Download the specified binary distribution package",
	Long: `
Allow simple download of the tarball/zip file for a specific version of a project. OS and
architecture are detected automatically, but can be overridden with the --os and --arch flags.
`,
	RunE: downloadPackage,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.ArgAliases = []string{"org", "project", "version"}
	downloadCmd.Args = cobra.ExactArgs(3)
	downloadCmd.Flags().StringP("arch", "a", runtime.GOARCH, "Force architecture (Default: current)")
	downloadCmd.Flags().StringP("os", "o", runtime.GOOS, "Force operating system (Default: current)")
	downloadCmd.Flags().StringP("path", "p", ".", "Path to store the downloaded package")
}

func downloadPackage(cmd *cobra.Command, args []string) error {

	org := args[0]
	project := args[1]
	version := args[2]

	os, _ := cmd.Flags().GetString("os")
	arch, _ := cmd.Flags().GetString("arch")
	path, _ := cmd.Flags().GetString("path")


	client := githubconnector.NewConnectorClient()
	path, err := client.DownloadPackage(org, project, version, os, arch, path)

	if err != nil {
		return errors.Wrap(err, "downloading package")
	}

	fmt.Println(path)
	return nil

}
