package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/chanzuckerberg/happy/shared/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ECSComputeBackend struct {
	Backend  *Backend
	SecretId string
}

type container struct {
	host          string
	container     string
	arn           string
	taskID        string
	launchType    string
	containerName string
}

type TaskInfo struct {
	TaskId     string `header:"Task ID"`
	StartedAt  string `header:"Started"`
	LastStatus string `header:"Status"`
}

const (
	AwsLogsGroup        = "awslogs-group"
	AwsLogsStreamPrefix = "awslogs-stream-prefix"
	AwsLogsRegion       = "awslogs-region"
)

type AWSLogConfiguration struct {
	GroupName   string
	StreamNames []string
}

func NewECSComputeBackend(ctx context.Context, secretId string, b *Backend) (*ECSComputeBackend, error) {
	return &ECSComputeBackend{
		Backend:  b,
		SecretId: secretId,
	}, nil
}

func (b *ECSComputeBackend) GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error) {
	out, err := b.Backend.secretsclient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &b.SecretId,
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not get integration secret at %s", b.SecretId)
	}

	secret := &config.IntegrationSecret{}
	err = json.Unmarshal([]byte(*out.SecretString), secret)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not json parse integraiton secret")
	}
	return secret, out.ARN, nil
}

func (b *ECSComputeBackend) GetParam(ctx context.Context, name string) (string, error) {
	log.Debugf("reading aws ssm parameter at %s", name)

	out, err := b.Backend.ssmclient.GetParameter(
		ctx,
		&ssm.GetParameterInput{Name: aws.String(name)},
	)
	if err != nil {
		return "", errors.Wrap(err, "could not get parameter")
	}

	return *out.Parameter.Value, nil
}

func (b *ECSComputeBackend) WriteParam(
	ctx context.Context,
	name string,
	val string,
) error {
	_, err := b.Backend.ssmclient.PutParameter(ctx, &ssm.PutParameterInput{
		Overwrite: aws.Bool(true),
		Name:      &name,
		Value:     &val,
		Type:      "String",
	})
	return errors.Wrapf(err, "could not write parameter to %s", name)
}

// GetLogGroupStreamsForStack gets all the log group and slice of log streams associated with a particular happy stack.
func (b *ECSComputeBackend) GetLogGroupStreamsForStack(ctx context.Context, stackName, serviceName string) (string, []string, error) {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "GetLogGroupStreamsForStack")

	stackTaskARNs, err := b.GetECSTasksForStackService(ctx, stackName, serviceName)
	if err != nil {
		return "", nil, errors.Wrapf(err, "error retrieving tasks for service '%s'", serviceName)
	}
	if len(stackTaskARNs) == 0 {
		return "", nil, errors.Errorf("no tasks associated with service %s. did you spell the service name correctly?", serviceName)
	}

	tasks, err := b.GetTaskDetails(ctx, stackTaskARNs)
	if err != nil {
		return "", nil, errors.Wrapf(err, "error getting task details for %+v", stackTaskARNs)
	}

	logConfigs, err := b.getAWSLogConfigsFromTasks(ctx, tasks...)
	if err != nil {
		return "", nil, err
	}

	return logConfigs.GroupName, logConfigs.StreamNames, nil
}

func (b *ECSComputeBackend) PrintLogs(ctx context.Context, stackName, serviceName, containerName string, opts ...util.PrintOption) error {
	// TODO: Add support for EKS log printing by container
	logGroup, logStreams, err := b.GetLogGroupStreamsForStack(ctx, stackName, serviceName)
	if err != nil {
		return err
	}

	opts = append([]util.PrintOption{util.WithPaginator(
		util.NewCloudWatchPaginator(cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:   &logGroup,
			LogStreamNames: logStreams,
		}, b.Backend.cwlFilterLogEventsAPIClient),
	)}, opts...)

	p := util.MakeComputeLogPrinter(ctx, opts...)
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "PrintLogs")
	err = p.Print(ctx)
	if err != nil {
		return errors.Wrap(err, "error printing logs")
	}

	expression := `fields @timestamp, @message
	| sort @timestamp desc
	| limit 20`

	logReference := util.LogReference{
		LinkOptions: util.LinkOptions{
			Region:       b.Backend.GetAWSRegion(),
			LaunchType:   util.LaunchTypeK8S,
			AWSAccountID: b.Backend.GetAWSAccountID(),
		},
		Expression:   expression,
		LogGroupName: logGroup,
	}

	return b.Backend.DisplayCloudWatchInsightsLink(ctx, logReference)
}

// RunTask runs an arbitrary task that is not necessarily associated with a service.
func (b *ECSComputeBackend) RunTask(ctx context.Context, taskDefArn string, launchType util.LaunchType) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), taskDefArn)

	log.Infof("running task %s, launch type %s", taskDefArn, launchType)
	out, err := b.Backend.ecsclient.RunTask(ctx, &ecs.RunTaskInput{
		Cluster:              &b.Backend.integrationSecret.ClusterArn,
		LaunchType:           ecstypes.LaunchType(launchType.String()),
		NetworkConfiguration: b.getNetworkConfig(),
		TaskDefinition:       &taskDefArn,
	})
	if err != nil {
		return errors.Wrapf(err, "could not run task %s", taskDefArn)
	}
	if len(out.Failures) > 0 {
		for _, failure := range out.Failures {
			log.Errorf("cannot start container %s because of %s", *failure.Arn, *failure.Reason)
		}
		return errors.Errorf("failed to launch task %s", taskDefArn)
	}
	if len(out.Tasks) == 0 {
		return errors.New("could not run task, not found")
	}

	tasks := []string{}
	for _, task := range out.Tasks {
		tasks = append(tasks, *task.TaskArn)
	}

	log.Infof("waiting for %+v to finish", tasks)
	err = b.waitForTasksToStop(ctx, tasks)
	if err != nil {
		return errors.Wrap(err, "error waiting for tasks")
	}

	log.Infof("task %s finished. printing logs from task", taskDefArn)
	logGroup, logStreams, err := b.getLogGroupStreamsForTasks(ctx, tasks)
	if err != nil {
		return err
	}

	since := util.GetStartTime(ctx).UnixMilli()
	p := util.MakeComputeLogPrinter(ctx,
		util.WithPaginator(
			util.NewCloudWatchPaginator(cloudwatchlogs.FilterLogEventsInput{
				LogGroupName:   &logGroup,
				LogStreamNames: logStreams,
				StartTime:      &since,
			}, b.Backend.cwlFilterLogEventsAPIClient),
		))
	return p.Print(ctx)
}

// GetLogGroupStreamsForTasks is just like GetLogGroupStreamsForService, except it gets all the log group and log
// streams associated with a one-off task, such as a migration or deletion task.
func (b *ECSComputeBackend) getLogGroupStreamsForTasks(ctx context.Context, taskARNs []string) (string, []string, error) {
	tasks, err := b.waitForTasksToStart(ctx, taskARNs)
	if err != nil {
		return "", nil, err
	}

	logConfigs, err := b.getAWSLogConfigsFromTasks(ctx, tasks...)
	if err != nil {
		return "", nil, err
	}

	return logConfigs.GroupName, logConfigs.StreamNames, nil
}

// getAWSLogConfigsFromTasks attempts to get all the log groups and log streams associated with the passed in
// tasks. It makes an assumption that all the tasks passed in are a part of the same happy service and therefore
// share a log group name. It only works for tasks that use the "awslog" driver.
func (b *ECSComputeBackend) getAWSLogConfigsFromTasks(ctx context.Context, tasks ...ecstypes.Task) (*AWSLogConfiguration, error) {
	logConfigs := &AWSLogConfiguration{
		GroupName:   "",
		StreamNames: []string{},
	}

	for _, task := range tasks {
		taskConfig, err := b.getAWSLogConfigsFromTask(ctx, task)
		if err != nil {
			return nil, err
		}
		logConfigs.StreamNames = append(logConfigs.StreamNames, taskConfig.StreamNames...)
		if logConfigs.GroupName == "" {
			logConfigs.GroupName = taskConfig.GroupName
		} else if logConfigs.GroupName != taskConfig.GroupName {
			return nil, errors.Errorf("expected the log groups to be the same. got %s and %s", logConfigs.GroupName, taskConfig.GroupName)
		}
	}
	return logConfigs, nil
}

// getAWSLogConfigsFromTask grabs all the log groups and log streams for all containers for a particular task.
func (b *ECSComputeBackend) getAWSLogConfigsFromTask(ctx context.Context, task ecstypes.Task) (*AWSLogConfiguration, error) {
	logConfigs := &AWSLogConfiguration{
		StreamNames: []string{},
	}
	tdef, err := b.Backend.ecsclient.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(*task.TaskDefinitionArn),
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not describe task definition")
	}

	for _, container := range tdef.TaskDefinition.ContainerDefinitions {
		// grab the awslogs logs-group and stream prefix from the task definition setting
		// ignore containers that don't have this or are not set to use 'awslogs' driver
		if container.LogConfiguration.LogDriver != "awslogs" {
			continue
		}
		var (
			awsLogStreamPrefix string
			ok                 bool
		)

		logConfigs.GroupName, ok = container.LogConfiguration.Options[AwsLogsGroup]
		if !ok {
			continue
		}

		awsLogStreamPrefix, ok = container.LogConfiguration.Options[AwsLogsStreamPrefix]
		// some tasks won't have a log stream prefix (a misconfiguration of their task definition)
		// TODO: let's try to always have a prefix, otherwise, we run into this issue:
		// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html
		// 	"If you don't specify a prefix with this option, then the log stream is named after the container ID
		// 	that's assigned by the Docker daemon on the container instance. Because it's difficult to trace logs
		// 	back to the container that sent them with just the Docker container ID (which is only available on
		// 	the container instance), we recommend that you specify a prefix with this option."
		if !ok {
			streams, err := b.Backend.cwlGetLogEventsAPIClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
				LogGroupName: &logConfigs.GroupName,
				OrderBy:      "LastEventTime",
				Descending:   aws.Bool(true),
			})
			if err != nil {
				return nil, errors.Wrapf(err, "unable to get log streams from log group %s", logConfigs.GroupName)
			}

			for _, stream := range streams.LogStreams {
				logConfigs.StreamNames = append(logConfigs.StreamNames, *stream.LogStreamName)
			}
		} else {
			taskID, err := b.getTaskID(*task.TaskArn)
			if err != nil {
				return nil, errors.Wrap(err, "unable to determine a task id")
			}
			// from the docs:
			// "Use the awslogs-stream-prefix option to associate a log stream with the specified prefix, the container name, and the ID
			// of the Amazon ECS task that the container belongs to. If you specify a prefix with this option, then the log stream takes the following format.
			// prefix-name/container-name/ecs-task-id"
			logConfigs.StreamNames = append(logConfigs.StreamNames, fmt.Sprintf("%s/%s/%s", awsLogStreamPrefix, *container.Name, taskID))
		}
	}

	return logConfigs, nil
}

// waitForTasksToStart waits for a set of tasks to no longer be in "PROVISIONING", "PENDING", or "ACTIVATING" states.
func (b *ECSComputeBackend) waitForTasksToStart(ctx context.Context, taskARNs []string) ([]ecstypes.Task, error) {
	tasks, err := util.IntervalWithTimeout(func() ([]ecstypes.Task, error) {
		tasks, err := b.GetTaskDetails(ctx, taskARNs)
		if err != nil {
			return nil, err
		}

		for _, task := range tasks {
			switch *task.LastStatus {
			case "PROVISIONING", "PENDING", "ACTIVATING":
				return nil, errors.Errorf("a task is not ready. %s still in %s", *task.TaskArn, *task.LastStatus)
			}
		}

		return tasks, nil
	}, 1*time.Second, 1*time.Minute)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		return nil, errors.New("unable to discover a task, impossible to stream the logs")
	}
	if tasks == nil || len(*tasks) == 0 {
		return nil, errors.Errorf("no matching tasks for task definition %+v", *tasks)
	}
	return *tasks, nil
}

func (b *ECSComputeBackend) waitForTasksToStop(ctx context.Context, taskARNs []string) error {
	descTask := &ecs.DescribeTasksInput{
		Cluster: &b.Backend.integrationSecret.ClusterArn,
		Tasks:   taskARNs,
	}
	err := b.Backend.taskStoppedWaiter.Wait(ctx, descTask, 600*time.Second)
	if err != nil {
		return errors.Wrap(err, "err waiting for tasks to stop")
	}

	tasks, err := b.Backend.ecsclient.DescribeTasks(ctx, descTask)
	if err != nil {
		return errors.Wrap(err, "could not describe tasks")
	}

	var failures error
	for _, failure := range tasks.Failures {
		failures = multierror.Append(failures, errors.Errorf("error running task (%s) with status (%s) and reason (%s)", *failure.Arn, *failure.Detail, *failure.Reason))
	}
	return failures
}

func (b *ECSComputeBackend) Shell(ctx context.Context, stackName, service, containerName, shellCommand string) error {
	clusterArn := b.Backend.Conf().GetClusterArn()

	serviceName := b.getEcsServiceName(stackName, service)
	ecsClient := b.Backend.GetECSClient()

	listTaskInput := &ecs.ListTasksInput{
		Cluster:     aws.String(clusterArn),
		ServiceName: aws.String(serviceName),
	}

	listTaskOutput, err := ecsClient.ListTasks(ctx, listTaskInput)
	if err != nil {
		return errors.Wrap(err, "error listing ecs tasks")
	}

	log.Println("Found tasks: ")
	tablePrinter := util.NewTablePrinter()

	describeTaskInput := &ecs.DescribeTasksInput{
		Cluster: aws.String(clusterArn),
		Tasks:   listTaskOutput.TaskArns,
	}

	describeTaskOutput, err := ecsClient.DescribeTasks(ctx, describeTaskInput)
	if err != nil {
		return errors.Wrap(err, "error describing ecs tasks")
	}

	containerMap := make(map[string]string)
	var containers []container

	for _, task := range describeTaskOutput.Tasks {
		taskArnSlice := strings.Split(*task.TaskArn, "/")
		taskID := taskArnSlice[len(taskArnSlice)-1]

		startedAt := "-"

		host := ""
		if task.ContainerInstanceArn != nil {
			host = *task.ContainerInstanceArn
		}

		if task.StartedAt != nil {
			startedAt = task.StartedAt.Format(time.RFC3339)
			containers = append(containers, container{
				host:          host,
				container:     *task.Containers[0].RuntimeId,
				arn:           *task.TaskArn,
				taskID:        taskID,
				launchType:    string(task.LaunchType),
				containerName: *task.Containers[0].Name,
			})
		}
		containerMap[*task.TaskArn] = host
		tablePrinter.AddRow(TaskInfo{TaskId: taskID, StartedAt: startedAt, LastStatus: *task.LastStatus})
	}

	tablePrinter.Flush()
	// FIXME: we make the assumption of only one container in many places. need consistency
	// TODO: only support ECS exec-command and NOT SSH
	for _, container := range containers {
		// This approach works for both Fargate and EC2 tasks
		awsProfile := b.Backend.GetAWSProfile()
		log.Infof("Connecting to %s:%s\n", container.taskID, container.containerName)
		// TODO: use the Go SDK and don't shell out
		//       see https://github.com/tedsmitt/ecsgo/blob/c1509097047a2d037577b128dcda4a35e23462fd/internal/pkg/internal.go#L196
		awsArgs := []string{"aws", "--profile", awsProfile, "ecs", "execute-command", "--cluster", clusterArn, "--container", container.containerName, "--command", "/bin/bash", "--interactive", "--task", container.taskID}

		awsCmd, err := b.Backend.executor.LookPath("aws")
		if err != nil {
			return errors.Wrap(err, "failed to locate the AWS cli")
		}

		cmd := &exec.Cmd{
			Path:   awsCmd,
			Args:   awsArgs,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		log.Println(cmd)
		if err := b.Backend.executor.Run(cmd); err != nil {
			return errors.Wrap(err, "failed to execute")
		}
	}

	return nil
}

func (b *ECSComputeBackend) getNetworkConfig() *ecstypes.NetworkConfiguration {
	privateSubnets := b.Backend.integrationSecret.PrivateSubnets
	privateSubnetsPt := []string{}
	for _, subnet := range privateSubnets {
		subnetValue := subnet
		privateSubnetsPt = append(privateSubnetsPt, subnetValue)
	}
	securityGroups := b.Backend.integrationSecret.SecurityGroups
	securityGroupsPt := []string{}
	for _, sg := range securityGroups {
		sgValue := sg
		securityGroupsPt = append(securityGroupsPt, sgValue)
	}

	awsvpcConfiguration := &ecstypes.AwsVpcConfiguration{
		AssignPublicIp: ecstypes.AssignPublicIpDisabled,
		SecurityGroups: securityGroupsPt,
		Subnets:        privateSubnetsPt,
	}
	networkConfig := &ecstypes.NetworkConfiguration{
		AwsvpcConfiguration: awsvpcConfiguration,
	}
	return networkConfig
}

func (b *ECSComputeBackend) getTaskID(taskARN string) (string, error) {
	resourceArn, err := arn.Parse(taskARN)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse task ARN: '%s'", taskARN)
	}

	segments := strings.Split(resourceArn.Resource, "/")
	if len(segments) < 3 {
		return "", errors.Errorf("incomplete task ARN: '%s'", taskARN)
	}
	return segments[len(segments)-1], nil
}

func (b *ECSComputeBackend) getEcsServiceName(stackName string, serviceName string) string {
	return fmt.Sprintf("%s-%s", stackName, serviceName)
}

func (b *ECSComputeBackend) GetEvents(ctx context.Context, stackName string, services []string) error {
	if len(services) == 0 {
		return nil
	}

	clusterArn := b.Backend.Conf().GetClusterArn()

	ecsClient := b.Backend.GetECSClient()

	ecsServices := make([]string, 0)
	for _, service := range services {
		ecsService := b.getEcsServiceName(stackName, service)
		ecsServices = append(ecsServices, ecsService)
	}

	describeServicesInput := &ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterArn),
		Services: ecsServices,
	}

	describeServicesOutput, err := ecsClient.DescribeServices(ctx, describeServicesInput)
	if err != nil {
		return errors.Wrap(err, "cannot describe services:")
	}

	for _, service := range describeServicesOutput.Services {
		incomplete := make([]string, 0)
		for _, deploy := range service.Deployments {
			if deploy.RolloutState != ecstypes.DeploymentRolloutStateCompleted {
				incomplete = append(incomplete, string(deploy.RolloutState))
			}
		}
		if len(incomplete) == 0 {
			continue
		}

		log.Println()
		log.Infof("Incomplete deployment of service %s / Current status %v:", *service.ServiceName, incomplete)

		deregistered := 0
		for index := range service.Events {
			event := service.Events[len(service.Events)-1-index]
			eventTime := event.CreatedAt
			if time.Since(*eventTime) > (time.Hour) {
				continue
			}

			message := regexp.MustCompile(`^\(service ([^ ]+)\)`).ReplaceAllString(*event.Message, "$1")
			message = regexp.MustCompile(`\(([^ ]+) .*?\)`).ReplaceAllString(message, "$1")
			message = regexp.MustCompile(`:.*`).ReplaceAllString(message, "$1")
			if strings.Contains(message, "deregistered") {
				deregistered++
			}

			log.Infof("  %s %s", eventTime.Format(time.RFC3339), message)
		}
		if deregistered > 3 {
			log.Println()
			log.Println("Many \"deregistered\" events - please check to see whether your service is crashing:")
			serviceName := strings.Replace(*service.ServiceName, fmt.Sprintf("%s-", stackName), "", 1)
			log.Infof("  happy --env %s logs %s %s", b.Backend.Conf().GetEnv(), stackName, serviceName)
		}
	}
	return nil
}

func isStackECSService(happyServiceName, happyStackName string, ecsService ecstypes.Service) bool {
	if strings.HasSuffix(*ecsService.ServiceName, happyServiceName) &&
		strings.HasPrefix(*ecsService.ServiceName, happyStackName) {
		return true
	}
	return false
}

// GetECSServicesForStackService returns the ECS services that are associated with a happy stack and service.
// The filter is based on the name of the stack and the service name provided in the docker-compose file.
func (b *ECSComputeBackend) GetECSServicesForStackService(ctx context.Context, stackName, serviceName string) ([]ecstypes.Service, error) {
	clusterARN := b.Backend.integrationSecret.ClusterArn

	var maxResults int32 = 10
	var services []ecstypes.Service
	var nextToken *string = nil
	for {
		ls, err := b.Backend.ecsclient.ListServices(ctx, &ecs.ListServicesInput{
			Cluster:    &clusterARN,
			MaxResults: &maxResults,
			NextToken:  nextToken,
		})
		if err != nil {
			break
		}
		if len(ls.ServiceArns) == 0 {
			break
		}
		ds, err := b.Backend.ecsclient.DescribeServices(ctx, &ecs.DescribeServicesInput{
			Cluster:  &clusterARN,
			Services: ls.ServiceArns,
		})
		if err != nil {
			return nil, errors.Wrap(err, "unable to describe ECS services for stack")
		}
		services = append(services, ds.Services...)
		if ls.NextToken == nil {
			break
		}
		nextToken = ls.NextToken
	}

	// TODO: right now, happy has no control over what these ECS services are called
	// but a convention has started where the stack name is a part of the service name
	// and so is the docker-compose service name. Usually, its of the form <stackname>-<docker-compose-service-name>.
	stackServNames := []ecstypes.Service{}
	for _, s := range services {
		if isStackECSService(serviceName, stackName, s) {
			stackServNames = append(stackServNames, s)
		}
	}

	return stackServNames, nil
}

// GetECSTasksForStackService returns the task ARNs associated with a particular happy stack and service.
func (b *ECSComputeBackend) GetECSTasksForStackService(ctx context.Context, stackName, serviceName string) ([]string, error) {
	stackServNames, err := b.GetECSServicesForStackService(ctx, stackName, serviceName)
	if err != nil {
		return nil, err
	}

	clusterARN := b.Backend.integrationSecret.ClusterArn
	stackTaskARNs := []string{}
	for _, s := range stackServNames {
		lt, err := b.Backend.ecsclient.ListTasks(ctx, &ecs.ListTasksInput{
			Cluster:     &clusterARN,
			ServiceName: s.ServiceName,
		})

		if err != nil {
			return nil, errors.Wrapf(err, "unable to list ECS tasks for stack %s", *s.ServiceName)
		}

		stackTaskARNs = append(stackTaskARNs, lt.TaskArns...)
	}

	log.Debugf("found task ARNs associated with %s-%s: %+v", stackName, serviceName, stackTaskARNs)
	return stackTaskARNs, nil
}

// GetTaskDefinitions returns the task definitions for a slice of task ARNs.
func (b *ECSComputeBackend) GetTaskDefinitions(ctx context.Context, taskArns []string) ([]ecstypes.TaskDefinition, error) {
	tasks, err := b.GetTaskDetails(ctx, taskArns)
	if err != nil {
		return []ecstypes.TaskDefinition{}, errors.Wrap(err, "cannot describe ECS tasks")
	}
	taskDefinitions := []ecstypes.TaskDefinition{}
	for _, task := range tasks {
		taskDefResult, err := b.Backend.ecsclient.DescribeTaskDefinition(
			ctx,
			&ecs.DescribeTaskDefinitionInput{TaskDefinition: task.TaskDefinitionArn},
		)
		if err != nil {
			return []ecstypes.TaskDefinition{}, errors.Wrap(err, "cannot retrieve a task definition")
		}
		taskDefinitions = append(taskDefinitions, *taskDefResult.TaskDefinition)
	}
	return taskDefinitions, nil
}

// GetTaskDetails is a helper function to return the described task objects associated with a slice
// of task ARNs.
func (b *ECSComputeBackend) GetTaskDetails(ctx context.Context, taskArns []string) ([]ecstypes.Task, error) {
	clusterARN := b.Backend.integrationSecret.ClusterArn
	tasksResult, err := b.Backend.ecsclient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   taskArns,
	})
	if err != nil {
		return []ecstypes.Task{}, errors.Wrap(err, "could not describe tasks")
	}
	return tasksResult.Tasks, nil
}

func (b *ECSComputeBackend) Describe(ctx context.Context, stackName string, serviceName string) (interfaces.StackServiceDescription, error) {
	params := make(map[string]string)
	params["cluster_arn"] = b.Backend.integrationSecret.ClusterArn
	params["service_name"] = b.getEcsServiceName(stackName, serviceName)
	params["integration_secret_id"] = b.SecretId
	description := interfaces.StackServiceDescription{
		Compute: "ECS",
		Params:  params,
	}
	return description, nil
}

func (b *ECSComputeBackend) GetResources(ctx context.Context, stackName string) ([]util.ManagedResource, error) {
	return []util.ManagedResource{}, nil
}
