package cmd

import (
	"fmt"
	"strings"

	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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

		tfeUrl := b.Conf().GetTfeUrl()
		tfeOrg := b.Conf().GetTfeOrg()

		workspaceRepo := workspace_repo.NewWorkspaceRepo(tfeUrl, tfeOrg)
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
		tablePrinter.AddRow("  Environment Workspace", fmt.Sprintf("%s/app/%s/workspaces/env-%s", tfeUrl, tfeOrg, bootstrapConfig.Env))
		tablePrinter.AddRow("  Stack Workspace", fmt.Sprintf("%s/app/%s/workspaces/%s-%s", tfeUrl, tfeOrg, bootstrapConfig.Env, stackName))

		tablePrinter.AddRow("AWS", "")
		tablePrinter.AddRow("  Account ID", fmt.Sprintf("[%s]", b.GetAWSAccountID()))
		tablePrinter.AddRow("  Region", b.GetAWSRegion())
		tablePrinter.AddRow("  Profile", b.GetAWSProfile())

		linkOptions := util.LinkOptions{
			Region:               b.GetAWSRegion(),
			IntegrationSecretARN: *b.GetIntegrationSecretArn(),
		}

		consoleUrl, err := util.Arn2ConsoleLink(linkOptions, b.Conf().ClusterArn)
		if err != nil {
			return errors.Errorf("error creating an AWS console link for ARN '%s'", b.Conf().ClusterArn)
		}

		tablePrinter.AddRow("ECS Cluster", consoleUrl)
		tablePrinter.AddRow("  ARN", b.Conf().ClusterArn)

		consoleUrl, err = util.Arn2ConsoleLink(linkOptions, *b.GetIntegrationSecretArn())
		if err != nil {
			return errors.Errorf("error creating an AWS console link for ARN '%s'", *b.GetIntegrationSecretArn())
		}
		tablePrinter.AddRow("Integration secret", consoleUrl)
		tablePrinter.AddRow("  ARN", *b.GetIntegrationSecretArn())

		for _, serviceName := range happyConfig.GetServices() {
			serviceName = fmt.Sprintf("%s-%s", stackName, serviceName)
			service, err := b.DescribeService(ctx, &serviceName)
			if err != nil {
				return errors.Errorf("error retrieving service details for service '%s'", serviceName)
			}
			consoleUrl, err := util.Arn2ConsoleLink(linkOptions, *service.ServiceArn)
			if err != nil {
				return errors.Errorf("error creating an AWS console link for ARN '%s'", *service.ServiceArn)
			}
			tablePrinter.AddRow("Service", consoleUrl)
			tablePrinter.AddRow("  Name", *service.ServiceName)
			tablePrinter.AddRow("  Launch Type", string(service.LaunchType))
			tablePrinter.AddRow("  Status", *service.Status)
			tablePrinter.AddRow("  Task Definition ARN", *service.TaskDefinition)
			tablePrinter.AddRow("    Desired Count", fmt.Sprintf("[%d]", service.DesiredCount))
			tablePrinter.AddRow("    Pending Count", fmt.Sprintf("[%d]", service.PendingCount))
			tablePrinter.AddRow("    Running Count", fmt.Sprintf("[%d]", service.RunningCount))

			taskArns, err := b.GetTasks(ctx, &serviceName)
			if err != nil {
				return errors.Wrapf(err, "error retrieving tasks for service '%s'", serviceName)
			}
			taskDefinitions, err := b.GetTaskDefinitions(ctx, taskArns)
			if err != nil {
				return errors.Wrapf(err, "error retrieving task definition for tasks '%v'", taskArns)
			}
			taskDefinitionMap := map[string]ecstypes.TaskDefinition{}
			for _, taskDefinition := range taskDefinitions {
				taskDefinitionMap[*taskDefinition.TaskDefinitionArn] = taskDefinition
			}

			tasks, err := b.GetTaskDetails(ctx, taskArns)
			if err != nil {
				return errors.Wrapf(err, "error retrieving task details for tasks '%s'", taskArns)
			}

			taskMap := map[string]ecstypes.Task{}
			for _, task := range tasks {
				taskMap[*task.TaskArn] = task
			}

			for _, taskArn := range taskArns {
				consoleUrl, err := util.Arn2ConsoleLink(linkOptions, taskArn)
				if err != nil {
					return errors.Wrapf(err, "error creating an AWS console link for ARN '%s'", taskArn)
				}
				tablePrinter.AddRow("  Task", consoleUrl)
				task := taskMap[taskArn]
				taskDefinition := taskDefinitionMap[*task.TaskDefinitionArn]

				arnSegments := strings.Split(taskArn, "/")
				if len(arnSegments) < 3 {
					continue
				}
				taskId := arnSegments[len(arnSegments)-1]
				tablePrinter.AddRow("    ARN", taskArn)
				tablePrinter.AddRow("    Status", *task.LastStatus)
				tablePrinter.AddRow("    Containers")
				for _, containerDefinition := range taskDefinition.ContainerDefinitions {
					tablePrinter.AddRow("      Name", *containerDefinition.Name)
					tablePrinter.AddRow("      Image", *containerDefinition.Image)

					logStreamPrefix := containerDefinition.LogConfiguration.Options[backend.AwsLogsStreamPrefix]
					logGroup := containerDefinition.LogConfiguration.Options[backend.AwsLogsGroup]
					logRegion := containerDefinition.LogConfiguration.Options[backend.AwsLogsRegion]
					containerName := *containerDefinition.Name

					consoleLink, err := util.Log2ConsoleLink(util.LinkOptions{Region: logRegion}, logGroup, logStreamPrefix, containerName, taskId)
					if err != nil {
						return errors.Errorf("unable to construct a cloudwatch link for container '%s'", containerName)
					}

					tablePrinter.AddRow("      Logs", consoleLink)
				}

			}
		}

		tablePrinter.Print()
		return nil
	},
}
