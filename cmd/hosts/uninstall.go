package hosts

import (
	"github.com/chanzuckerberg/happy/pkg/hostname_manager"
	"github.com/spf13/cobra"
)

var unInstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove compose DNS entries",
	Long:  "Remove compose DNS entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostManager := hostname_manager.NewHostNameManager(hostsFile, nil)
		return hostManager.UnInstall()
	},
}
