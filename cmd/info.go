package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
	config.ConfigureCmdWithBootstrapConfig(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:          "info",
	Short:        "info [stack]",
	Long:         "Get information on resources allocated for the environment '{env}' and the stack, if specified",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.RangeArgs(0, 1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		stackName := ""
		if len(args) == 1 {
			stackName = args[0]
		}

		bootstrapConfig, err := config.NewBootstrapConfig(cmd)
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

		url := b.Conf().GetTfeUrl()
		org := b.Conf().GetTfeOrg()

		headings := []string{"Resource", "Value"}
		tablePrinter := util.NewTablePrinter(headings)

		tablePrinter.AddRow("TFE Environment Workspaces", fmt.Sprintf("%s/app/%s/workspaces", url, org))

		tablePrinter.AddRow("AWS Region", b.GetAWSRegion())
		tablePrinter.AddRow("AWS Profile", b.GetAWSProfile())

		consoleUrl, err := arn2ConsoleLink(b, b.Conf().ClusterArn)
		tablePrinter.AddRow("ECS Cluster", consoleUrl)
		tablePrinter.AddRow("ECS Cluster ARN", b.Conf().ClusterArn)

		consoleUrl, err = arn2ConsoleLink(b, b.GetIntegrationSecret().GetSecretArn())
		tablePrinter.AddRow("Integration secret", consoleUrl)
		tablePrinter.AddRow("Integration secret ARN", b.GetIntegrationSecret().GetSecretArn())

		tablePrinter.AddRow("Environment", bootstrapConfig.Env)
		if len(stackName) > 0 {
			tablePrinter.AddRow("Stack", stackName)
			tablePrinter.AddRow("TFE Workspace", fmt.Sprintf("%s/app/%s/workspaces/%s-%s", url, org, bootstrapConfig.Env, stackName))
		}

		tablePrinter.Print()

		return err
	},
}

func arn2ConsoleLink(b *backend.Backend, unparsedArn string) (string, error) {
	resourceArn, err := arn.Parse(unparsedArn)
	if err != nil {
		return "", errors.Wrapf(err, "Invalid ARN: %s", unparsedArn)
	}

	region := b.GetAWSRegion()

	switch resourceArn.Service {
	case "ecs":

		resourceParts := strings.Split(resourceArn.Resource, "/")
		return fmt.Sprintf("https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services", region, region, resourceParts[1]), nil

	case "secretsmanager":
		//arn:aws:secretsmanager:us-west-2:626314663667:secret:happy/env-ep-config-2E78OI
		//https://us-west-2.console.aws.amazon.com/secretsmanager/home?region=us-west-2#!/secret?name=happy%2Fenv-ep-config

		secretArn := strings.ReplaceAll(url.QueryEscape(b.Conf().HappyConfig.GetSecretArn()), "%", "%%")

		return fmt.Sprintf("https://%s.console.aws.amazon.com/secretsmanager/home?region=%s#!/secret?name=%s", region, region, secretArn), nil
	}

	return "", errors.Errorf("service %s is not supported", unparsedArn)

}
