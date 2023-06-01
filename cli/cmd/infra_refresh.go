package cmd

import (
	"github.com/chanzuckerberg/happy/shared/config"
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

		if len(happyConfig.GetData().StackDefaults) == 0 {
			err = hclManager.Ingest(ctx)
			if err != nil {
				return errors.Wrap(err, "unable to ingest HCL code")
			}
			happyConfig, err = config.GetHappyConfigForCmd(cmd)
			if err != nil {
				return err
			}
			hclManager.WithHappyConfig(happyConfig)
		}

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

		err = hclManager.Generate(ctx)
		return errors.Wrap(err, "unable to generate HCL code")
	},
}
