package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
	config.ConfigureCmdWithBootstrapConfig(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "for test",
	Long:  "for testing",
	RunE:  runCmd,
}

func runCmd(cmd *cobra.Command, args []string) error {
	bootstrapConfig, err := config.NewBootstrapConfig()
	if err != nil {
		return err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return err
	}

	clusterArn := happyConfig.ClusterArn()
	privateSubnets := happyConfig.PrivateSubnets()
	securityGroups := happyConfig.SecurityGroups()

	fmt.Println("This is the cluster ARN: ", clusterArn)
	fmt.Println("This is the private subnets: ", privateSubnets)
	fmt.Println("This is the security groups: ", securityGroups)

	return nil
}
