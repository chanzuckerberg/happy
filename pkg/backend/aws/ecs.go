package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	AwsLogsGroup        = "awslogs-group"
	AwsLogsStreamPrefix = "awslogs-stream-prefix"
	AwsLogsRegion       = "awslogs-region"
)

func (b *Backend) DescribeService(ctx context.Context, serviceName *string) (*ecstypes.Service, error) {
	clusterARN := b.integrationSecret.ClusterArn
	out, err := b.ecsclient.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Services: []string{*serviceName},
		Cluster:  &clusterARN,
		Include:  []ecstypes.ServiceField{},
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot describe an ECS service")
	}
	if len(out.Services) == 0 {
		return nil, errors.New("ECS service was not found")
	}

	return &out.Services[0], nil
}

func (b *Backend) GetServiceTasks(ctx context.Context, serviceName *string) ([]string, error) {
	clusterARN := b.integrationSecret.ClusterArn
	out, err := b.ecsclient.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:     &clusterARN,
		ServiceName: serviceName,
	})
	if err != nil {
		var snfe *ecstypes.ServiceNotFoundException
		if errors.As(err, &snfe) {
			return []string{}, nil
		}
		return []string{}, errors.Wrap(err, "cannot retrieve ECS tasks")
	}
	return out.TaskArns, nil
}

func (b *Backend) GetTaskDefinitions(ctx context.Context, taskArns []string) ([]ecstypes.TaskDefinition, error) {
	clusterARN := b.integrationSecret.ClusterArn

	tasksResult, err := b.ecsclient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   taskArns,
	})

	if err != nil {
		return []ecstypes.TaskDefinition{}, errors.Wrap(err, "cannot describe ECS tasks")
	}
	taskDefinitions := []ecstypes.TaskDefinition{}
	for _, task := range tasksResult.Tasks {
		taskDefResult, err := b.ecsclient.DescribeTaskDefinition(
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

func (b *Backend) GetTaskDetails(ctx context.Context, taskArns []string) ([]ecstypes.Task, error) {
	clusterARN := b.integrationSecret.ClusterArn
	tasksResult, err := b.ecsclient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   taskArns,
	})
	if err != nil {
		return []ecstypes.Task{}, errors.Wrap(err, "could not describe tasks")
	}
	return tasksResult.Tasks, nil
}

func (b *Backend) RunTask(
	ctx context.Context,
	taskDefArn string,
	launchType config.LaunchType,
) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), taskDefArn)
	clusterARN := b.integrationSecret.ClusterArn
	networkConfig := b.getNetworkConfig()

	log.Infof("running task %s", taskDefArn)
	out, err := b.ecsclient.RunTask(ctx, &ecs.RunTaskInput{
		Cluster:              &clusterARN,
		LaunchType:           ecstypes.LaunchType(launchType.String()),
		NetworkConfiguration: networkConfig,
		TaskDefinition:       &taskDefArn,
	})
	if err != nil {
		return errors.Wrapf(err, "could not run task %s", taskDefArn)
	}

	if len(out.Tasks) == 0 {
		return errors.New("could not run task, not found")
	}

	tasks := []string{}
	for _, task := range out.Tasks {
		tasks = append(tasks, *task.TaskArn)
	}

	descTasks := &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   tasks,
	}

	err = b.waitForTasks(ctx, descTasks)
	if err != nil {
		return errors.Wrap(err, "error waiting for tasks")
	}
	log.Infof("task %s finished. printing logs from task", taskDefArn)
	// log the tasks after they are done
	return b.getLogEventsForTask(ctx, taskDefArn, descTasks, util.MakeLogPrinter())
}

func (ab *Backend) waitForTasks(ctx context.Context, input *ecs.DescribeTasksInput) error {
	err := ab.taskStoppedWaiter.Wait(ctx, input, 600*time.Second)
	if err != nil {
		return errors.Wrap(err, "err waiting for tasks to stop")
	}

	// now get their status
	tasks, err := ab.ecsclient.DescribeTasks(ctx, input)
	if err != nil {
		return errors.Wrap(err, "could not describe tasks")
	}

	var failures error
	for _, failure := range tasks.Failures {
		failures = multierror.Append(failures, errors.Errorf("error running task (%s) with status (%s) and reason (%s)", *failure.Arn, *failure.Detail, *failure.Reason))
	}
	return failures
}

func (ab *Backend) getNetworkConfig() *ecstypes.NetworkConfiguration {
	privateSubnets := ab.integrationSecret.PrivateSubnets
	privateSubnetsPt := []string{}
	for _, subnet := range privateSubnets {
		subnetValue := subnet
		privateSubnetsPt = append(privateSubnetsPt, subnetValue)
	}
	securityGroups := ab.integrationSecret.SecurityGroups
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

func (ab *Backend) Logs(ctx context.Context, stackName string, serviceName string, since string) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "Logs")
	taskArns, err := ab.GetServiceTasks(ctx, &serviceName)
	if err != nil {
		return errors.Wrapf(err, "error retrieving tasks for service '%s'", serviceName)
	}

	tasks, err := ab.waitAndDescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &ab.integrationSecret.ClusterArn,
		Tasks:   taskArns,
	})
	if err != nil {
		return errors.Wrapf(err, "error waiting for tasks %s, %+v", serviceName, taskArns)
	}

	logConfigs, err := ab.getAWSLogConfigsFromTasks(ctx, tasks.Tasks...)
	if err != nil {
		return err
	}

	params := cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:   &logConfigs.GroupName,
		LogStreamNames: logConfigs.StreamNames,
	}
	if since != "" {
		duration, err := time.ParseDuration(since)
		if err != nil {
			return errors.Wrapf(err, "unable to parse the 'since' param %s", since)
		}
		params.StartTime = aws.Int64(getStartTime(ctx).Add(-duration).UnixMilli())
	}

	log.Debugf("Following logs: group=%s, stream=%+v", logConfigs.GroupName, logConfigs.StreamNames)
	return ab.GetLogs(ctx, &params, util.MakeLogPrinter())
}

func intervalWithTimeout[K any](f func() (*K, error), tick time.Duration, timeout time.Duration) (*K, error) {
	timeoutChan := time.After(timeout)
	tickChan := time.NewTicker(tick)

	for {
		select {
		case <-timeoutChan:
			return nil, errors.New("timed out")
		case <-tickChan.C:
			out, err := f()
			if err == nil {
				return out, nil
			}
			log.Infof("trying again: %s", err)
		}
	}
}

func (ab *Backend) waitAndDescribeTasks(ctx context.Context, input *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	tasks, err := intervalWithTimeout(func() (*ecs.DescribeTasksOutput, error) {
		t, err := ab.ecsclient.DescribeTasks(ctx, input)
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, errors.New("task isn't ready")
		}

		for _, task := range (*t).Tasks {
			switch *task.LastStatus {
			case "PROVISIONING", "PENDING", "ACTIVATING":
				return nil, errors.Errorf("a task is not ready. %s still in %s", *task.TaskArn, *task.LastStatus)
			}
		}

		// if all the tasks are ready, return them
		return t, nil
	}, 1*time.Second, 1*time.Minute)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		return nil, errors.New("unable to discover a task, impossible to stream the logs")
	}
	if tasks == nil || len(tasks.Tasks) == 0 {
		return nil, errors.Errorf("no matching tasks for task definition %+v", input.Tasks)
	}
	return tasks, nil
}

func getStartTime(ctx context.Context) time.Time {
	// This is the value the task was started, we don't want logs before this
	// time.
	cmdStartTime, ok := ctx.Value(util.CmdStartContextKey).(time.Time)
	if !ok {
		log.Debugf("didn't get a cmd start time. using now")
		cmdStartTime = time.Now()
	}
	return cmdStartTime
}

func (ab *Backend) getLogEventsForTask(
	ctx context.Context,
	taskDefARN string,
	input *ecs.DescribeTasksInput,
	filterLogs FilterLogsFunc,
) error {
	tasks, err := ab.waitAndDescribeTasks(ctx, input)
	if err != nil {
		return err
	}

	logConfigs, err := ab.getAWSLogConfigsFromTasks(ctx, tasks.Tasks...)
	if err != nil {
		return err
	}
	
	return ab.GetLogs(
		ctx,
		&cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:   &logConfigs.GroupName,
			LogStreamNames: logConfigs.StreamNames,
			StartTime:      aws.Int64(getStartTime(ctx).UnixMilli()),
		},
		filterLogs,
	)
}

func (ab *Backend) getTaskID(taskARN string) (string, error) {
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

type AWSLogConfiguration struct {
	GroupName   string
	StreamNames []string
}

func (ab *Backend) getAWSLogConfigsFromTask(ctx context.Context, task ecstypes.Task) (*AWSLogConfiguration, error) {
	logConfigs := &AWSLogConfiguration{
		StreamNames: []string{},
	}
	tdef, err := ab.ecsclient.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
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
			streams, err := ab.cwlGetLogEventsAPIClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
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
			taskID, err := ab.getTaskID(*task.TaskArn)
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

// getAWSLogConfigsFromTasks attempts to get all the log groups and log streams associated with the passed in
// tasks. It makes an assumption that all the tasks passed in are a part of the same service and therefore
// share a log group name. It only works for tasks that use the "awslog" driver. It also only grabs the latest
// log stream if the task definition of a task is without an "awslogs-stream-prefix".
func (ab *Backend) getAWSLogConfigsFromTasks(ctx context.Context, tasks ...ecstypes.Task) (*AWSLogConfiguration, error) {
	logConfigs := &AWSLogConfiguration{
		GroupName:   "",
		StreamNames: []string{},
	}

	for _, task := range tasks {
		taskConfig, err := ab.getAWSLogConfigsFromTask(ctx, task)
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
