package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(availableVersionCmd)

	versionCmd.AddCommand(lockHappyVersionCmd)
	lockHappyVersionCmd.Flags().String("version", "", "Specify a version of Happy to lock in .happy/version.lock file. Default to current CLI version.")
}

var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Version of Happy",
	Long:         "Returns the current version of the happy cli",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		v := util.GetVersion().String()
		fmt.Fprintln(cmd.OutOrStdout(), v)
		return nil
	},
}

var availableVersionCmd = &cobra.Command{
	Use:          "available-version",
	Short:        "Latest available version of Happy",
	Long:         "Returns the latest available version of Happy",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := GetLatestAvailableVersion(cmd)
		if err != nil {
			fmt.Fprint(cmd.Parent().ErrOrStderr(), err)
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), v)
		return nil
	},
}

var lockHappyVersionCmd = &cobra.Command{
	Use:          "lock",
	Short:        "Create a .happy/version.lock file",
	Long:         "Create a .happy/version.lock file in project root to specify which version of Happy should be used with this project. This will overwrite any existing version.lock file.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		versionFile, version, err := CreateHappyVersionFile(cmd)
		if err != nil {
			fmt.Fprint(cmd.Parent().ErrOrStderr(), err)
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Created %s locking Happy to version %s\n", versionFile, version)
		return nil
	},
}

func GetLatestAvailableVersion(cmd *cobra.Command) (*util.Release, error) {
	happyConfig, err := config.GetHappyConfigForCmd(cmd)
	if err != nil {
		return nil, err
	}

	api := hapi.MakeApiClient(happyConfig)
	result := model.HealthResponse{}
	err = api.GetParsed("/health", "", &result)
	if err != nil {
		return nil, err
	}

	return &util.Release{
		Version: result.Version,
		GitSha:  result.GitSha,
	}, nil
}

func IsHappyOutdated(cmd *cobra.Command) (bool, *util.Release, *util.Release, error) {
	cliVersion := util.GetVersion()
	latestAvailableVersion, err := GetLatestAvailableVersion(cmd)
	if err != nil {
		return false, cliVersion, latestAvailableVersion, err
	}

	return !cliVersion.Equal(latestAvailableVersion), cliVersion, latestAvailableVersion, nil
}

func WarnIfHappyOutdated(cmd *cobra.Command) {

	outdated, cliVersion, latestAvailableVersion, err := IsHappyOutdated(cmd)

	if err != nil {
		log.Errorf("Error checking for latest available version number: %v", err)
		return
	}

	if outdated {
		log.Warnf("This copy of Happy CLI is not the latest available. CLI version: %s  Latest available version: %s\n", cliVersion.Version, latestAvailableVersion.Version)
		log.Warn("To update on Mac, run:  brew upgrade happy")
	}

}

func CreateHappyVersionFile(cmd *cobra.Command) (string, string, error) {
	happyConfig, err := config.GetHappyConfigForCmd(cmd)
	if err != nil {
		return "", "", err
	}

	currentVersion := util.GetVersion()
	projectRoot := happyConfig.GetProjectRoot()

	requestedVersion, _ := cmd.Flags().GetString("version")

	if requestedVersion == "" {
		requestedVersion = currentVersion.Version
	}

	versionFile := config.NewHappyVersionLockFile(projectRoot)

	err = versionFile.SetVersion(requestedVersion)
	if err != nil {
		return "", "", err
	}

	err = versionFile.Save()
	if err != nil {
		return "", "", err
	}

	versionFilePath, err := versionFile.GetPath()
	if err != nil {
		return "", "", err
	}

	version, err := versionFile.GetVersion()
	if err != nil {
		return "", "", err
	}

	return versionFilePath, version, nil
}
