/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/user"
	"path"
	"runtime"

	"github.com/chanzuckerberg/happy/hvm/installer"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a version of Happy",
	Long:  `Install a version of Happy to ~/.happy/versions/ and set it as the current version.`,
	Run:   installPackage,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.ArgAliases = []string{"versionTag"}
	installCmd.Args = cobra.ExactArgs(1)
	installCmd.Flags().StringP("arch", "a", "", "Force architecture (Default: current)")
	installCmd.Flags().StringP("os", "o", "", "Force operating system (Default: current)")
	installCmd.Flags().StringP("path", "p", ".", "Path to store the downloaded package")

}

func installPackage(cmd *cobra.Command, args []string) {

	versionTag := args[0]
	os := runtime.GOOS
	arch := runtime.GOARCH

	if cmd.Flags().Changed("os") {
		os = cmd.Flag("os").Value.String()
	}

	if cmd.Flags().Changed("arch") {
		arch = cmd.Flag("arch").Value.String()
	}

	user, err := user.Current()

	if err != nil {
		fmt.Println("Error getting current user information", err)
		return
	}

	home := user.HomeDir

	versionsPath := path.Join(home, ".czi", "versions", "happy", versionTag)

	if cmd.Flags().Changed("path") {
		versionsPath = cmd.Flag("path").Value.String()
	}

	err = installer.InstallPackage(versionTag, os, arch, versionsPath)

	if err != nil {
		fmt.Println("Error installing package", err)
		return
	}

}
