package cmd

import (
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infraCmd)
	config.ConfigureCmdWithBootstrapConfig(infraCmd)
}

var infraCmd = &cobra.Command{
	Use:          "infra",
	Short:        "Infra commands",
	Long:         "Execute infra commands in environment '{env}'",
	SilenceUsage: false,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
			checklist.AwsInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Info("Please specify a subcommand. See --help for more information.")
		return nil
	},
}
