package hosts

import (
	"github.com/chanzuckerberg/happy/cli/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/cli/pkg/hostname_manager"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/spf13/cobra"
)

func init() {
	config.ConfigureCmdWithBootstrapConfig(installCmd)
}

var installCmd = &cobra.Command{
	Use:          "install",
	Short:        "Install compose DNS entries",
	Long:         "Install compose DNS entries",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		bootstrapConfig, err := config.NewBootstrapConfig(cmd)
		if err != nil {
			return err
		}
		happyConfig, err := config.NewHappyConfig(bootstrapConfig)
		if err != nil {
			return err
		}

		buildConfig := artifact_builder.NewBuilderConfig().WithBootstrap(bootstrapConfig).WithHappyConfig(happyConfig)
		containers, err := buildConfig.GetContainers(ctx)
		if err != nil {
			return err
		}

		hostManager := hostname_manager.NewHostNameManager(hostsFile, containers)
		return hostManager.Install()

	},
}
