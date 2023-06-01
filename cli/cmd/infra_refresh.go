package cmd

import (
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	infraCmd.AddCommand(infraValidateCmd)
	config.ConfigureCmdWithBootstrapConfig(infraRefreshCmd)
}

var infraRefreshCmd = &cobra.Command{
	Use:          "refresh",
	Short:        "Refresh Happy Stack HCL code",
	Long:         "Refresh Happy Stack HCL code for all environments",
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

		logrus.Debug("Refreshing HCL code")
		err = hclManager.Ingest(ctx)
		if err != nil {
			return errors.Wrap(err, "Unable to ingest HCL code")
		}
		err = hclManager.Generate(ctx)
		return errors.Wrap(err, "Unable to generate HCL code")
	},
}
