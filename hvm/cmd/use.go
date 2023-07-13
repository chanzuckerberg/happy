package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Download locked version of Happy",
	Long: `If a Happy version lock file exists in the project config directory, download the specified version.
Otherwise download the latest available version of Happy.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use called")
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	useCmd.Flags().String("bin-path", "/usr/local/bin", "Path to store the happy binary")
	useCmd.Flags().String("arch", "", "Download for a specific architecture (Default: current)")
	useCmd.Flags().String("os", "", "Download for a specific operating system (Default: current)")

}
