package cmd

import (
	"fmt"
	"os"

	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "current version of the happy cli",
	Long:         "returns the current version of the happy cli",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		v := util.GetVersion().String()
		fmt.Fprintln(os.Stdout, v)
		return nil
	},
}
