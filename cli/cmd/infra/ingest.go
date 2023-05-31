package infra

import (
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	infraCmd.AddCommand(infraIngestCmd)
	config.ConfigureCmdWithBootstrapConfig(infraIngestCmd)
}

var infraIngestCmd = &cobra.Command{
	Use:          "ingest",
	Short:        "Ingest Happy Stack HCL code",
	Long:         "Ingest Happy Stack HCL code from environment '{env}' into happy config",
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

		logrus.Debug("Ingesting HCL code")
		return hclManager.Ingest(ctx)
	},
}
