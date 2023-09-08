/*
 */
package cmd

import (
	"os"
	"path"
	"runtime"

	"github.com/chanzuckerberg/happy/hvm/installer"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [org] [project] [version]",
	Short: "Install a version of a project",
	Long:  `Install a version of a project to ~/.happy/versions/ and set it as the current version.`,
	RunE:  installPackage,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.ArgAliases = []string{"org", "project", "version"}
	installCmd.Args = cobra.ExactArgs(3)
	installCmd.Flags().StringP("arch", "a", runtime.GOARCH, "Force architecture (Default: current)")
	installCmd.Flags().StringP("os", "o", runtime.GOOS, "Force operating system (Default: current)")

}

func installPackage(cmd *cobra.Command, args []string) error {

	org := args[0]
	project := args[1]
	version := args[2]

	opsys := cmd.Flag("os").Value.String()
	arch := cmd.Flag("arch").Value.String()

	home, err := os.UserHomeDir()

	if err != nil {
		return errors.Wrap(err, "getting current user home directory")
	}

	versionsPath := path.Join(home, ".czi", "versions", org, project, version)

	return installer.InstallPackage(cmd.Context(), org, project, version, opsys, arch, versionsPath)

}
