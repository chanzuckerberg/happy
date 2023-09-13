package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock [org] [project] [version]",
	Short: "Lock a specific version of a requirement in the current project",
	Long:  `Lock a specific version of a requirement in the current project. This will create a .happy/version.lock file`,
	RunE:  setLock,
}

func init() {
	rootCmd.AddCommand(lockCmd)
	lockCmd.ArgAliases = []string{"org", "project", "version"}
	lockCmd.Args = cobra.ExactArgs(3)
}

func setLock(cmd *cobra.Command, args []string) error {

	org := args[0]
	project := args[1]
	version := args[2]

	happyConfig, err := config.GetHappyConfigForCmd(cmd)
	if err != nil {
		return errors.Wrap(err, "getting happy config")
	}

	projectRoot := happyConfig.GetProjectRoot()

	lockfile, err := config.NewHappyVersionLockFile(projectRoot)

	if err != nil {
		return errors.Wrap(err, "creating default version lockfile")
	}

	if config.DoesHappyVersionLockFileExist(projectRoot) {
		lockfile, err = config.LoadHappyVersionLockFile(projectRoot)
		if err != nil {
			return errors.Wrap(err, "loading version lockfile")
		}
	}

	lockSlug := fmt.Sprintf("%s/%s", org, project)
	lockfile.Require[lockSlug] = version

	err = lockfile.Save()
	if err != nil {
		return errors.Wrap(err, "saving version lockfile")
	}

	fmt.Printf("Locked %s to %s\n", lockSlug, version)

	return nil

}
