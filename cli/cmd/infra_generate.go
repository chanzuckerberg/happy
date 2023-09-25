package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/shared/composemanager"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
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
		composeManager := composemanager.NewComposeManager().WithHappyConfig(happyConfig)

		if !force {
			if !happyConfig.GetData().FeatureFlags.EnableUnifiedConfig {
				if diagnostics.IsInteractiveContext(ctx) {
					proceed := false
					prompt := &survey.Confirm{Message: "Currently, stack settings are managed in terraform code. Are you sure you want to overwrite them?"}
					err = survey.AskOne(prompt, &proceed)
					if err != nil {
						return errors.Wrapf(err, "failed to ask for confirmation")
					}

					if !proceed {
						return err
					}
				}
			}
		}

		logrus.Debug("Generating HCL code")
		err = hclManager.Generate(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to generate HCL code")
		}
		logrus.Debug("Generating docker-compose file")
		return errors.Wrap(composeManager.Manage(ctx), "unable to generate docker-compose file")
	},
}
