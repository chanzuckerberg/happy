package cmd

import (
	"fmt"
	"strings"

	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
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
	Use:          "set-lock",
	Short:        "Create a .happy/version.lock file",
	Long:         "Create a .happy/version.lock file in project root to specify which version of Happy should be used with this project. This will overwrite any existing version.lock file.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := CreateHappyVersionLockfileHandler(cmd)
		if err != nil {
			log.Debug(cmd.Parent().ErrOrStderr(), err)
			return err
		}

		return nil
	},
}

func GetLatestAvailableVersion(cmd *cobra.Command) (*util.Release, error) {
	happyClient, err := makeHappyClient(cmd, sliceName, "", []string{}, false)
	if err != nil {
		return nil, err
	}

	api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend)
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
		log.Debug("Error getting latest available version. Will fail silently.")
		return false, cliVersion, cliVersion, nil // Lie.
	}

	cliVersionStr := normalizeVersion(cliVersion.Version)
	latestAvailableVersionStr := normalizeVersion(latestAvailableVersion.Version)
	outdated := semver.Compare(cliVersionStr, latestAvailableVersionStr) < 0

	return outdated, cliVersion, latestAvailableVersion, nil
}

func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	}
	return fmt.Sprintf("v%s", version)
}

func WarnIfHappyOutdated(cmd *cobra.Command) {
	outdated, cliVersion, latestAvailableVersion, err := IsHappyOutdated(cmd)

	if err != nil {
		log.Errorf("Error checking for latest available version number: %v", err)
		return
	}

	if outdated {
		log.Warnf("This copy of Happy CLI is not the latest available. CLI version: %s  Latest available version: %s", cliVersion.Version, latestAvailableVersion.Version)
		log.Warn("To update on Mac, run:  brew upgrade happy")
	}
}

func CreateHappyVersionLockfileHandler(cmd *cobra.Command) error {
	happyConfig, err := config.GetHappyConfigForCmd(cmd)
	if err != nil {
		return err
	}

	currentVersion := util.GetVersion()
	projectRoot := happyConfig.GetProjectRoot()

	requestedVersion, _ := cmd.Flags().GetString("version")

	if requestedVersion == "" {
		requestedVersion = currentVersion.Version
	}

	path, version, err := createHappyVersionLockFile(projectRoot, requestedVersion)

	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created %s locking Happy to version %s\n", path, version)

	return nil
}

func createHappyVersionLockFile(projectRoot string, requestedVersion string) (string, string, error) {
	versionFile, err := config.NewHappyVersionLockFile(projectRoot)
	if err != nil {
		return "", "", err
	}

	versionFile.Require["chanzuckerberg/happy"] = requestedVersion
	err = versionFile.Save()
	if err != nil {
		return "", "", err
	}

	return versionFile.VersionSpecPath, requestedVersion, nil
}

func VerifyHappyIsLockedVersion(cmd *cobra.Command) (bool, string, string, error) {
	happyConfig, err := config.GetHappyConfigForCmd(cmd)
	if err != nil {
		return false, "", "", err
	}

	projectRoot := happyConfig.GetProjectRoot()

	/*
		For backward compatibility reasons, if the .happy/version.lock file does not exist,
		we will simply return true. We are essentially saying that an unlocked version is the same
		as when an version is locked and matched.
	*/

	if !config.DoesHappyVersionLockFileExist(projectRoot) {
		return true, "", "", nil
	}

	happyVersionLock, err := config.LoadHappyVersionLockFile(projectRoot)
	if err != nil {
		return false, "", "", err
	}

	if util.GetVersion().Version != happyVersionLock.Require["chanzuckerberg/happy"] {
		return false, util.GetVersion().Version, happyVersionLock.Require["chanzuckerberg/happy"], nil
	}

	return true, util.GetVersion().Version, happyVersionLock.Require["chanzuckerberg/happy"], nil
}

func CheckLockedHappyVersion(cmd *cobra.Command) error {
	excludeLockedVersionCheckCmds := map[string]interface{}{
		"version":           nil,
		"set-lock":          nil,
		"available-version": nil,
	}

	log.Debugf("Current command: %s\n", cmd.CalledAs())
	if _, present := excludeLockedVersionCheckCmds[cmd.CalledAs()]; present {
		log.Debug("Skipping locked version check")
		return nil
	}

	versionMatch, cliVersion, lockedVersion, err := VerifyHappyIsLockedVersion(cmd)
	if err != nil {
		// This is generally going to be because we are outside of a project root.
		log.Debugf("Unable to verify locked Happy version: %s", err)
		return nil
	}

	// version will be undefined for local builds of the CLI and we don't want it to break here when testing locally
	if !versionMatch && cliVersion != "undefined" {
		return errors.Errorf("installed Happy version (%s) does not match locked version in .happy/version.lock (%s)", cliVersion, lockedVersion)
	}

	return nil
}
