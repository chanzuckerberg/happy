package cmd

import (
	"github.com/chanzuckerberg/happy/cli/pkg/config_manager"
	"github.com/chanzuckerberg/happy/shared/composemanager"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/hclmanager"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	config.ConfigureCmdWithBootstrapConfig(bootstrapCmd)
	bootstrapCmd.Flags().BoolVar(&force, "force", false, "Ignore the already-exists errors")
}

var bootstrapCmd = &cobra.Command{
	Use:          "bootstrap",
	Short:        "Bootstrap the happy repo",
	Long:         "Configure the repo to be used with happy",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		checklist := util.NewValidationCheckList()
		return util.ValidateEnvironment(cmd.Context(),
			checklist.TerraformInstalled,
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		bootstrapConfig, err := config.NewSimpleBootstrap(cmd)
		if err == nil && !force {
			return errors.New("this repo is already bootstrapped")
		}

		happyConfig, err := config_manager.CreateHappyConfig(ctx, bootstrapConfig)
		if err != nil {
			return errors.Wrap(err, "unable to create a new happy config")
		}

		hclManager := hclmanager.NewHclManager().WithHappyConfig(happyConfig)
		composeManager := composemanager.NewComposeManager().WithHappyConfig(happyConfig)

		log.Debug("Generating HCL code")
		for envName := range happyConfig.GetData().Environments {
			happyConfig.SetEnv(envName)
			err = hclManager.Generate(ctx)
			if err != nil {
				log.Errorf("unable to generate HCL code for env %s: %s", envName, err.Error())
			}
		}
		log.Debug("Generating docker-compose file")
		return errors.Wrap(composeManager.Manage(ctx), "unable to generate docker-compose file")
	},
}
