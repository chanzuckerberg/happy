package cmd

import (
	"fmt"
	"os"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "for test",
	Long:  "for testing",
	RunE:  runCmd,
}

func runCmd(cmd *cobra.Command, args []string) error {

	env := "rdev"

	happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
	if !ok {
		return errors.New("please set env var HAPPY_CONFIG_PATH")
	}

	happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
	if err != nil {
		return err
	}

	clusterArn, err := happyConfig.ClusterArn()
	if err != nil {
		return err
	}
	privateSubnets, err := happyConfig.PrivateSubnets()
	if err != nil {
		return err
	}
	securityGroups, err := happyConfig.SecurityGroups()
	if err != nil {
		return err
	}
	fmt.Println("This is the cluster ARN: ", clusterArn)
	fmt.Println("This is the private subnets: ", privateSubnets)
	fmt.Println("This is the security groups: ", securityGroups)

	return nil
}
