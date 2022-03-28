package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/cmd"
	"github.com/chanzuckerberg/happy/pkg/config"
	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
	config.ConfigureCmdWithBootstrapConfig(getCmd)
}

var getCmd = &cobra.Command{
	Use:          "get",
	Short:        "get stack",
	Long:         "Get a stack in environment '{env}'",
	SilenceUsage: true,
	PreRunE:      cmd.Validate(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		stackName := args[0]

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

		workspaceRepo := workspace_repo.NewWorkspaceRepo(url, org)
		stackSvc := stackservice.NewStackService().WithBackend(b).WithWorkspaceRepo(workspaceRepo)

		stacks, err := stackSvc.GetStacks(ctx)
		if err != nil {
			return err
		}

		stack, ok := stacks[stackName]
		if !ok {
			return errors.Errorf("stack '%s' not found in environment '%s'", stackName, happyConfig.GetEnv())
		}

		logrus.Infof("Retrieving stack '%s' from environment '%s'", stackName, happyConfig.GetEnv())

		headings := []string{"Name", "Owner", "Tags", "Status", "URLs"}
		tablePrinter := util.NewTablePrinter(headings)

		err = stack.Print(ctx, stackName, tablePrinter)
		if err != nil {
			logrus.Errorf("Error retrieving stack %s:  %s", stackName, err)
		}

		tablePrinter.Print()

		headings = []string{"Resource", "Value"}
		tablePrinter = util.NewTablePrinter(headings)

		tablePrinter.AddRow("Environment", bootstrapConfig.Env)
		tablePrinter.AddRow("TFE", "")
		tablePrinter.AddRow("  Environment Workspace", fmt.Sprintf("%s/app/%s/workspaces/env-%s", url, org, bootstrapConfig.Env))
		tablePrinter.AddRow("  Stack Workspace", fmt.Sprintf("%s/app/%s/workspaces/%s-%s", url, org, bootstrapConfig.Env, stackName))

		tablePrinter.AddRow("AWS", "")
		tablePrinter.AddRow("  Account ID", fmt.Sprintf("%s.", b.GetAWSAccountID()))
		tablePrinter.AddRow("  Region", b.GetAWSRegion())
		tablePrinter.AddRow("  Profile", b.GetAWSProfile())

		consoleUrl, err := arn2ConsoleLink(b, b.Conf().ClusterArn)
		tablePrinter.AddRow("ECS Cluster", consoleUrl)
		tablePrinter.AddRow("  ARN", b.Conf().ClusterArn)

		consoleUrl, err = arn2ConsoleLink(b, b.GetIntegrationSecret().GetSecretArn())
		tablePrinter.AddRow("Integration secret", consoleUrl)
		tablePrinter.AddRow("  ARN", b.GetIntegrationSecret().GetSecretArn())

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
		secretName := strings.ReplaceAll(url.QueryEscape(b.Conf().HappyConfig.GetSecretArn()), "%", "%%")
		return fmt.Sprintf("https://%s.console.aws.amazon.com/secretsmanager/home?region=%s#!/secret?name=%s", region, region, secretName), nil
	}

	return "", errors.Errorf("service %s is not supported", unparsedArn)
}
