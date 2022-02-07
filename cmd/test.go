package cmd

import (
	"github.com/chanzuckerberg/happy/pkg/config"
	log "github.com/sirupsen/logrus"
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

	log.Println("This is the cluster ARN: ", clusterArn)
	log.Println("This is the private subnets: ", privateSubnets)
	log.Println("This is the security groups: ", securityGroups)

	return nil
}
