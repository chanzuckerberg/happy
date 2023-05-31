package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(apiCmd)
	apiCmd.AddCommand(apiHealthCmd)
}

var apiCmd = &cobra.Command{
	Use:          "api",
	Short:        "meta-command for happy api",
	Long:         "Meta-command for checking on the api",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Println(cmd.Usage())
		return nil
	},
}

var apiHealthCmd = &cobra.Command{
	Use:          "ping",
	Short:        "ping health",
	Long:         "ping the health endpoint of happy api",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		api := hapi.MakeApiClient(happyConfig)
		result := model.HealthResponse{}
		err = api.GetParsed("/health", "", &result)
		if err != nil {
			return err
		}
		logrus.Infof("happy-api (%s%s) status: %s latest available version: %s", happyConfig.GetHappyApiConfig().BaseUrl, result.Route, result.Status, result.Version)

		return nil
	},
}
