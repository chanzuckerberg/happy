package cmd

import (
	"runtime"

	"github.com/chanzuckerberg/happy/hvm/githubconnector"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [versionTag]",
	Short: "Download the specified binary distribution package for Happy",
	Long: `
Allow simple download of the tarball/zip file for a specific version of Happy. OS and 
architecture are detected automatically, but can be overridden with the --os and --arch flags.
`,
	Run: downloadPackage,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.ArgAliases = []string{"versionTag"}
	downloadCmd.Args = cobra.ExactArgs(1)
	downloadCmd.Flags().StringP("arch", "a", "", "Force architecture (Default: current)")
	downloadCmd.Flags().StringP("os", "o", "", "Force operating system (Default: current)")
	downloadCmd.Flags().StringP("path", "p", ".", "Path to store the downloaded package")
}

func downloadPackage(cmd *cobra.Command, args []string) {
	versionTag := args[0]

	os := runtime.GOOS
	arch := runtime.GOARCH
	path := "."

	if cmd.Flags().Changed("os") {
		os = cmd.Flag("os").Value.String()
	}

	if cmd.Flags().Changed("arch") {
		arch = cmd.Flag("arch").Value.String()
	}

	if cmd.Flags().Changed("path") {
		arch = cmd.Flag("path").Value.String()
	}

	client := githubconnector.NewConnectorClient()
	client.DownloadPackage(versionTag, os, arch, path)

}
