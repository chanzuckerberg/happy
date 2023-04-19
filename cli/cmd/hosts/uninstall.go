package hosts

import (
	"github.com/chanzuckerberg/happy/cli/pkg/hostname_manager"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/spf13/cobra"
)

func init() {
	config.ConfigureCmdWithBootstrapConfig(unInstallCmd)
}

var unInstallCmd = &cobra.Command{
	Use:          "uninstall",
	Short:        "Remove compose DNS entries",
	Long:         "Remove compose DNS entries",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostManager := hostname_manager.NewHostNameManager(hostsFile, nil)
		return hostManager.UnInstall()
	},
}
