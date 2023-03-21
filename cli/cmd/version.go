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
	rootCmd.AddCommand(availableVersionCmd)
}

var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Version of Happy",
	Long:         "Returns the current version of the happy cli",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		v := util.GetVersion().String()
		fmt.Fprint(cmd.OutOrStdout(), v)
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

		fmt.Fprint(cmd.OutOrStdout(), v)
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
