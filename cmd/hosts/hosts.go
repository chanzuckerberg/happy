package hosts

import "github.com/spf13/cobra"

var hostsFile string

func NewHostsCommand() *cobra.Command {
	hostsCmd := &cobra.Command{
		Use:   "hosts",
		Short: "Commands to manage system hostsfile",
		Long:  "Commands to manage system hostsfile",
	}

	hostsCmd.AddCommand(installCmd)
	hostsCmd.AddCommand(unInstallCmd)

	installCmd.Flags().StringVar(&hostsFile, "hostsfile", "/etc/hosts", "Path to system hosts file (default is /etc/hosts)")
	unInstallCmd.Flags().StringVar(&hostsFile, "hostsfile", "/etc/hosts", "Path to system hosts file (default is /etc/hosts)")
}
