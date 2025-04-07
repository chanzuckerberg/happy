package cmd

import (
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/util/tf"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pin string

func init() {
	infraCmd.AddCommand(infraRefreshCmd)
	infraRefreshCmd.Flags().StringVar(&pin, "pin", "main", "Git tag to pin the happy stack module to")
	config.ConfigureCmdWithBootstrapConfig(infraRefreshCmd)
}

var infraRefreshCmd = &cobra.Command{
	Use:          "refresh",
	Short:        "Refresh Happy Stack HCL code",
	Long:         "Refresh Happy Stack HCL code in environment '{env}'",
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
						err = happyConfig.Save()
						if err != nil {
							return errors.Wrapf(err, "failed to save happy config")
						}
						hclManager.WithHappyConfig(happyConfig)
					}
				}
			}
		}

		err = hclManager.Ingest(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to ingest HCL code")
		}
		happyConfig, err = config.GetHappyConfigForCmd(cmd)
		if err != nil {
			return err
		}
		hclManager.WithHappyConfig(happyConfig)

		moduleSource := happyConfig.GetModuleSource()

		if len(moduleSource) == 0 {
			return errors.New("module source cannot be determined")
		}

		log.Debugf("module source: %s", moduleSource)

		gitUrl, path, _, err := tf.ParseModuleSource(moduleSource)
		if err != nil {
			return errors.Wrapf(err, "unable to parse module source: %s", moduleSource)
		}

		updatedModuleSource := tf.ComposeModuleSource(gitUrl, path, pin)

		if moduleSource != updatedModuleSource {
			log.Debugf("updated module source to: %s", updatedModuleSource)
			envConfig := happyConfig.GetData().Environments[happyConfig.GetEnv()]
			if envConfig.StackOverrides == nil {
				envConfig.StackOverrides = map[string]interface{}{}
			}

			envConfig.StackOverrides["source"] = moduleSource
			happyConfig.GetData().Environments[happyConfig.GetEnv()] = envConfig
			err = happyConfig.Save()
			if err != nil {
				return errors.Wrap(err, "unable to save happy config")
			}
		}

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

		err = hclManager.Generate(ctx)
		return errors.Wrap(err, "unable to generate HCL code")
	},
}
