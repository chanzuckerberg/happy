package cmd

import (
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	infraCmd.AddCommand(infraGenerateCmd)
	config.ConfigureCmdWithBootstrapConfig(infraGenerateCmd)
}

var infraGenerateCmd = &cobra.Command{
	Use:          "generate",
	Short:        "Generate Happy Stack HCL code",
	Long:         "Generate Happy Stack HCL code in environment '{env}'",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

		logrus.Debug("Generating HCL code")
		return hclManager.Generate(ctx)
	},
}
