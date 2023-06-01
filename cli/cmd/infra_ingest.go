package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	infraCmd.AddCommand(infraIngestCmd)
	config.ConfigureCmdWithBootstrapConfig(infraIngestCmd)
}

var infraIngestCmd = &cobra.Command{
	Use:          "ingest",
	Short:        "Ingest Happy Stack HCL code",
	Long:         "Ingest Happy Stack HCL code from all environments",
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

		if len(happyConfig.GetData().StackDefaults) > 0 {
			log.Info("happy config already has stack defaults, they will be overwritten")
		}

		hclManager := hclmanager.NewHclManager().WithHappyConfig(happyConfig)

		log.Debug("Ingesting HCL code")

		if !force {
			if diagnostics.IsInteractiveContext(ctx) {
				if happyConfig.GetData().FeatureFlags.EnableUnifiedConfig {
					proceed := false
					prompt := &survey.Confirm{Message: "Stack settings are managed in happy config, this will overwrite your existing stack defaults. Are you sure you want to proceed?"}
					err = survey.AskOne(prompt, &proceed)
					if err != nil {
						return errors.Wrapf(err, "failed to ask for confirmation")
					}

					if !proceed {
						return err
					}
				} else {
					proceed := false
					prompt := &survey.Confirm{Message: "Would you like to manage stack settings in happy config instead of terraform code?"}
					err = survey.AskOne(prompt, &proceed)
					if err != nil {
						return errors.Wrapf(err, "failed to ask for confirmation")
					}

					if proceed {
						happyConfig.GetData().FeatureFlags.EnableUnifiedConfig = true
						happyConfig.Save()
						hclManager.WithHappyConfig(happyConfig)
					}
				}
			}
		}

		return hclManager.Ingest(ctx)
	},
}
