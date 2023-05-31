package cmd

import (
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ciCmd)
	config.ConfigureCmdWithBootstrapConfig(ciCmd)
}

var ciCmd = &cobra.Command{
	Use:          "ci-role",
	Short:        "Get the CI role",
	Long:         "Print the happy environment's CI role ARN to be used in Github Actions",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		happyClient, err := makeHappyClient(cmd, sliceName, "", []string{tag}, createTag)
		if err != nil {
			return errors.Wrap(err, "unable to initialize the happy client")
		}
		ciRoles := happyClient.AWSBackend.Conf().IntegrationSecret.CIRoles
		if ciRoles == nil {
			return nil
		}

		logrus.Info("CI Roles:")
		for _, role := range *ciRoles {
			logrus.Infof("\t- %s\n", role.RoleARN)
		}
		return nil
	},
}
