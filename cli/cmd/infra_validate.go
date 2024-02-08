package cmd

import (
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	infraCmd.AddCommand(infraValidateCmd)
	config.ConfigureCmdWithBootstrapConfig(infraValidateCmd)
}

var infraValidateCmd = &cobra.Command{
	Use:          "validate",
	Short:        "Validate Happy Stack HCL code",
	Long:         "Validate Happy Stack HCL code for all environments",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
			checklist.AwsInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		happyConfig, err := config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}

		hclManager := hclmanager.NewHclManager().WithHappyConfig(happyConfig)

		logrus.Debug("Validating HCL code")
		return hclManager.Validate(ctx, hclmanager.ENV_ALL)
	},
}
