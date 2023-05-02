package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ciCmd)
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

		fmt.Println("CI Roles:")
		for _, role := range *ciRoles {
			fmt.Printf("\t- %s\n", role.RoleARN)
		}
		return nil
	},
}
