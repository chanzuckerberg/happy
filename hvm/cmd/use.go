package cmd

import (
	"fmt"

	linkmanager "github.com/chanzuckerberg/happy/hvm/linkManager"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Download locked version of Happy",
	Long: `If a Happy version lock file exists in the project config directory, download the specified version.
Otherwise download the latest available version of Happy.
	`,
	Run: useVersion,
}

func init() {
	rootCmd.AddCommand(useCmd)

}

func useVersion(cmd *cobra.Command, args []string) {
	versionTag := args[0]

	err := linkmanager.SetBinLink(versionTag)

	if err != nil {
		fmt.Println("Error setting bin link: ", err)
	}
}
