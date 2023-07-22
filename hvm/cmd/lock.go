package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the current version of happy in the current project",
	Long:  `Lock the current version of happy in the current project. This will create a .happy/version.lock file`,
	Run:   setLock,
}

func init() {
	rootCmd.AddCommand(lockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func setLock(cmd *cobra.Command, args []string) {
	fmt.Println("UNIMPLEMENTED: lock called")
}
