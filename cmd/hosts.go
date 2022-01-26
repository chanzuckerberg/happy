package cmd

import (
	"errors"
	"os"

	"github.com/chanzuckerberg/happy/pkg/artifact_builder"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/hostname_manager"
	"github.com/spf13/cobra"
)

var hostsFile string

func init() {
	rootCmd.AddCommand(hostsCmd)
	hostsCmd.AddCommand(installCmd)
	hostsCmd.AddCommand(unInstallCmd)
	installCmd.Flags().StringVar(&hostsFile, "hostsfile", "/etc/hosts", "Path to system hosts file (default is /etc/hosts)")
	unInstallCmd.Flags().StringVar(&hostsFile, "hostsfile", "/etc/hosts", "Path to system hosts file (default is /etc/hosts)")
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install compose DNS entries",
	Long:  "Install compose DNS entries",
	RunE: func(cmd *cobra.Command, args []string) error {

		env := "rdev"

		dockerComposeConfigPath, ok := os.LookupEnv("DOCKER_COMPOSE_CONFIG_PATH")
		if !ok {
			return errors.New("please set env var DOCKER_COMPOSE_CONFIG_PATH")
		}

		happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
		if !ok {
			return errors.New("please set env var HAPPY_CONFIG_PATH")
		}

		happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
		if err != nil {
			return err
		}

		composeEnv := ""
		if useComposeEnv {
			composeEnv = happyConfig.DefaultComposeEnv()
		}
		buildConfig := artifact_builder.NewBuilderConfig(dockerComposeConfigPath, composeEnv)

		containers := buildConfig.GetContainers()
		hostManager := hostname_manager.NewHostNameManager(hostsFile, containers)
		hostManager.Install()
		return nil

	},
}

var unInstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove compose DNS entries",
	Long:  "Remove compose DNS entries",
	RunE: func(cmd *cobra.Command, args []string) error {

		hostManager := hostname_manager.NewHostNameManager(hostsFile, nil)
		hostManager.UnInstall()
		return nil
	},
}

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Commands to manage system hostsfile",
	Long:  "Commands to manage system hostsfile",
}
