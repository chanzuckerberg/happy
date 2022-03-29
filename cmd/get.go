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
			return errors.Wrapf(err, "error retrieving stack '%s'", stackName)
		}

		tablePrinter.Print()

		headings = []string{"Resource", "Value"}
		tablePrinter = util.NewTablePrinter(headings)

		tablePrinter.AddRow("Environment", bootstrapConfig.Env)
		tablePrinter.AddRow("TFE", "")
		tablePrinter.AddRow("  Environment Workspace", fmt.Sprintf("%s/app/%s/workspaces/env-%s", url, org, bootstrapConfig.Env))
		tablePrinter.AddRow("  Stack Workspace", fmt.Sprintf("%s/app/%s/workspaces/%s-%s", url, org, bootstrapConfig.Env, stackName))

		tablePrinter.AddRow("AWS", "")
		tablePrinter.AddRow("  Account ID", fmt.Sprintf("[%s]", b.GetAWSAccountID()))
		tablePrinter.AddRow("  Region", b.GetAWSRegion())
		tablePrinter.AddRow("  Profile", b.GetAWSProfile())

		consoleUrl, err := arn2ConsoleLink(b, b.Conf().ClusterArn)
		if err != nil {
			return errors.Errorf("error creating an AWS console link for ARN '%s'", b.Conf().ClusterArn)
		}

		tablePrinter.AddRow("ECS Cluster", consoleUrl)
		tablePrinter.AddRow("  ARN", b.Conf().ClusterArn)

		consoleUrl, err = arn2ConsoleLink(b, b.GetIntegrationSecret().GetSecretArn())
		if err != nil {
			return errors.Errorf("error creating an AWS console link for ARN '%s'", b.GetIntegrationSecret().GetSecretArn())
		}
		tablePrinter.AddRow("Integration secret", consoleUrl)
		tablePrinter.AddRow("  ARN", b.GetIntegrationSecret().GetSecretArn())

		for _, serviceName := range happyConfig.GetServices() {
			serviceName = fmt.Sprintf("%s-%s", stackName, serviceName)
			service, err := b.DescribeService(ctx, &serviceName)
			if err != nil {
				return errors.Errorf("error retrieving service details for service '%s'", serviceName)
			}
			consoleUrl, err := arn2ConsoleLink(b, *service.ServiceArn)
			if err != nil {
				return errors.Errorf("error creating an AWS console link for ARN '%s'", *service.ServiceArn)
			}
			tablePrinter.AddRow("Service", consoleUrl)
			tablePrinter.AddRow("  Name", *service.ServiceName)
			tablePrinter.AddRow("  Launch Type", *service.LaunchType)
			tablePrinter.AddRow("  Status", *service.Status)
			tablePrinter.AddRow("  Task Definition ARN", *service.TaskDefinition)
			tablePrinter.AddRow("    Desired Count", fmt.Sprintf("[%d]", *service.DesiredCount))
			tablePrinter.AddRow("    Pending Count", fmt.Sprintf("[%d]", *service.PendingCount))
			tablePrinter.AddRow("    Running Count", fmt.Sprintf("[%d]", *service.RunningCount))

			taskArns, err := b.GetTasks(ctx, &serviceName)
			if err != nil {
				return errors.Errorf("error retrieving tasks for service '%s'", serviceName)
			}
			for _, taskArn := range taskArns {
				consoleUrl, err := arn2ConsoleLink(b, *taskArn)
				if err != nil {
					return errors.Errorf("error creating an AWS console link for ARN '%s'", *taskArn)
				}
				tablePrinter.AddRow("  Task", consoleUrl)

				taskDefinitions, err := b.GetTaskDefinitions(ctx, taskArn)
				if err != nil {
					return errors.Errorf("error retrieving task definition for task '%s'", *taskArn)
				}
				tasks, err := b.GetTaskDetails(ctx, taskArn)
				if err != nil {
					return errors.Errorf("error retrieving task details for task '%s'", *taskArn)
				}

				for taskIndex, taskDefinition := range taskDefinitions {
					task := tasks[taskIndex]
					arnSegments := strings.Split(*taskArn, "/")
					if len(arnSegments) < 3 {
						continue
					}
					taskId := arnSegments[len(arnSegments)-1]
					tablePrinter.AddRow("    ARN", *taskArn)
					tablePrinter.AddRow("    Status", *task.LastStatus)
					tablePrinter.AddRow("    Containers")
					for _, containerDefinition := range taskDefinition.ContainerDefinitions {
						tablePrinter.AddRow("      Name", *containerDefinition.Name)
						tablePrinter.AddRow("      Image", *containerDefinition.Image)

						logStreamPrefix := *containerDefinition.LogConfiguration.Options["awslogs-stream-prefix"]
						logGroup := *containerDefinition.LogConfiguration.Options["awslogs-group"]
						logRegion := *containerDefinition.LogConfiguration.Options["awslogs-region"]
						containerName := *containerDefinition.Name

						link := fmt.Sprintf("https://%s.console.aws.amazon.com/cloudwatch/home?region=%s#logEventViewer:group=%s;stream=%s/%s/%s", logRegion, logRegion, logGroup, logStreamPrefix, containerName, taskId)
						tablePrinter.AddRow("      Logs", link)
					}
				}
			}
		}

		tablePrinter.Print()
		return nil
	},
}

func arn2ConsoleLink(b *backend.Backend, unparsedArn string) (string, error) {
	resourceArn, err := arn.Parse(unparsedArn)
	if err != nil {
		return "", errors.Wrapf(err, "invalid ARN: %s", unparsedArn)
	}

	region := b.GetAWSRegion()

	switch resourceArn.Service {
	case "ecs":
		resourceParts := strings.Split(resourceArn.Resource, "/")
		if len(resourceParts) < 2 {
			return "", errors.Wrapf(err, "ARN is not supported: %s", unparsedArn)
		}
		resourceType := resourceParts[0]
		resourceName := resourceParts[1]
		resourceSubName := ""
		if len(resourceParts) > 2 {
			resourceSubName = resourceParts[2]
		}

		switch resourceType {
		case "cluster":
			return fmt.Sprintf("https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services", region, region, resourceName), nil
		case "task":
			return fmt.Sprintf("https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/tasks/%s/details", region, region, resourceName, resourceSubName), nil
		case "service":
			return fmt.Sprintf("https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services/%s/tasks", region, region, resourceName, resourceSubName), nil
		}
		return "", errors.Errorf("resource %s is not supported", resourceType)

	case "secretsmanager":
		resourceParts := strings.Split(resourceArn.Resource, ":")
		resourceType := resourceParts[0]
		switch resourceType {
		case "secret":
			secretName := strings.ReplaceAll(url.QueryEscape(b.Conf().HappyConfig.GetSecretArn()), "%", "%%")
			return fmt.Sprintf("https://%s.console.aws.amazon.com/secretsmanager/home?region=%s#!/secret?name=%s", region, region, secretName), nil
		}
		return "", errors.Errorf("resource %s is not supported", resourceType)
	}

	return "", errors.Errorf("service %s is not supported", unparsedArn)
}
