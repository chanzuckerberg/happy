package cmd

import (
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
	config.ConfigureCmdWithBootstrapConfig(testCmd)
}

var testCmd = &cobra.Command{
	Use:          "test",
	Short:        "for test",
	Long:         "for testing",
	SilenceUsage: true,
	RunE:         runCmd,
}

func runCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	bootstrapConfig, err := config.NewBootstrapConfig()
	if err != nil {
		return err
	}
	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	if err != nil {
		return err
	}

	b, err := backend.NewAWSBackend(ctx, happyConfig)
	if err != nil {
		return err
	}

	clusterArn := b.Conf().GetClusterArn()
	privateSubnets := b.Conf().GetPrivateSubnets()
	securityGroups := b.Conf().GetSecurityGroups()

	logrus.Info("This is the cluster ARN: ", clusterArn)
	logrus.Info("This is the private subnets: ", privateSubnets)
	logrus.Info("This is the security groups: ", securityGroups)

	return nil
}
