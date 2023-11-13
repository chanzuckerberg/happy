/*
HVM Version Commands
*/
package cmd

import (
	"fmt"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output version of HVM",
	Long:  `Output the current version of the HVM CLI`,
	RunE:  outputVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func outputVersion(cmd *cobra.Command, args []string) error {
	v := util.GetVersion().String()
	fmt.Fprintln(cmd.OutOrStdout(), v)
	return nil
}
