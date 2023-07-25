package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock [org] [project] [version]",
	Short: "Lock the current version of happy in the current project",
	Long:  `Lock the current version of happy in the current project. This will create a .happy/version.lock file`,
	Run:   setLock,
}

func init() {
	rootCmd.AddCommand(lockCmd)
	lockCmd.ArgAliases = []string{"org", "project", "version"}
	lockCmd.Args = cobra.ExactArgs(3)
}

func setLock(cmd *cobra.Command, args []string) {

	org := args[0]
	project := args[1]
	version := args[2]

	var lockfile *config.HappyVersionLockFile

	happyConfig, err := config.GetHappyConfigForCmd(cmd)
	if err != nil {
		fmt.Println("Error getting happy config: ", err)
		return
	}

	projectRoot := happyConfig.GetProjectRoot()

        lockfile, err := config.NewHappyVersionLockFile(projectRoot)
	if config.DoesHappyVersionLockFileExist(projectRoot) {
		lockfile, err = config.LoadHappyVersionLockFile(projectRoot)
		if err != nil {
			return errors.Wrap(err, "loading version lockfile")
	}

	if err != nil {
		fmt.Println("Error getting version lockfile: ", err)
		return
	}

	lockSlug := fmt.Sprintf("%s/%s", org, project)
	lockfile.Require[lockSlug] = version

	err = lockfile.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Locked %s to %s\n", lockSlug, version)

}
