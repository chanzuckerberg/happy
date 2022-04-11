package aws

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/chanzuckerberg/happy/pkg/config"
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

func (b *Backend) GetTasks(ctx context.Context, serviceName *string) ([]string, error) {
	clusterARN := b.integrationSecret.ClusterArn
	out, err := b.ecsclient.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:     &clusterARN,
		ServiceName: serviceName,
	})
	if err != nil {
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
	clusterARN := b.integrationSecret.ClusterArn
	networkConfig := b.getNetworkConfig()

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

	waitInput := &ecs.DescribeTasksInput{
		Cluster: &clusterARN,
		Tasks:   tasks,
	}

	// start reading logs asynchronously
	done := make(chan struct{})
	logserr := make(chan error, 1)

	go func() {
		logserr <- b.getLogEventsForTask(
			ctx,
			taskDefArn,
			waitInput,
			func(gleo *cloudwatchlogs.GetLogEventsOutput, err error) error {
				if err != nil {
					return err
				}
				select {
				case <-done:
					return Stop()
				default:
					for _, event := range gleo.Events {
						// TODO: better output here
						log.Info(*event.Message)
					}
					return nil
				}
			},
		)
	}()

	err = b.waitForTasks(ctx, waitInput)
	close(done)
	if err != nil {
		return errors.Wrap(err, "error waiting for tasks")
	}

	return <-logserr
}

func (ab *Backend) waitForTasks(ctx context.Context, input *ecs.DescribeTasksInput) error {
	err := ab.taskRunningWaiter.Wait(ctx, input, 600*time.Second)
	if err != nil {
		return errors.Wrap(err, "err waiting for tasks to start")
	}

	err = ab.taskStoppedWaiter.Wait(ctx, input, 600*time.Second)
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

func (ab *Backend) Logs(ctx context.Context, serviceName string, since string) error {
	endTime := time.Now()

	duration, err := time.ParseDuration(since)
	if err != nil {
		return errors.Wrapf(err, "invalid duration: '%s'", since)
	}
	startTime := endTime.Add(-duration)

	logGroup := ""
	logStreamName := ""

	// Get a list of task ARNs for a given service
	taskArns, err := ab.GetTasks(ctx, &serviceName)
	if err != nil {
		return errors.Wrapf(err, "error retrieving tasks for service '%s'", serviceName)
	}
	tasks, err := ab.GetTaskDetails(ctx, taskArns)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve task: '%s'", taskArns)
	}

	taskDefinitions, err := ab.GetTaskDefinitions(ctx, taskArns)
	if err != nil {
		return errors.Wrapf(err, "error retrieving task definition for task '%v'", taskArns)
	}
	taskDefinitionMap := map[string]ecstypes.TaskDefinition{}
	for _, taskDefinition := range taskDefinitions {
		taskDefinitionMap[*taskDefinition.TaskDefinitionArn] = taskDefinition
	}
	taskMap := map[string]ecstypes.Task{}
	for _, task := range tasks {
		taskMap[*task.TaskArn] = task
	}

	// For now container name is hardcoded (look across all containers)
	containerName := ""
	// Look for the first task definition that has a log group for a matching container name
	for _, taskArn := range taskArns {
		taskId, err := ab.getTaskId(taskArn)
		if err != nil {
			return errors.Wrapf(err, "invalid task ARN: '%s'", taskArn)
		}

		task := taskMap[taskArn]
		taskDefinition := taskDefinitionMap[*task.TaskDefinitionArn]
		logGroup, logStreamName, err = ab.getLogGroupAndStreamName(taskDefinition, taskId, containerName)
		if err != nil {
			log.Debugf("task definition %s does not have a log group: %s", *taskDefinition.TaskDefinitionArn, err.Error())
			continue
		}

		if len(logGroup) > 0 && len(logStreamName) > 0 {
			break
		}
	}

	if len(logGroup) == 0 {
		return errors.Errorf("unable to determine a log group for service '%s'", serviceName)
	}

	if len(logStreamName) == 0 {
		return errors.Errorf("unable to determine a log stream name for service '%s'", serviceName)
	}

	params := cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &logGroup,
		LogStreamName: &logStreamName,
		StartTime:     aws.Int64(startTime.UnixNano() / int64(time.Millisecond)),
		EndTime:       aws.Int64(endTime.UnixNano() / int64(time.Millisecond)),
	}

	log.Infof("Tailing logs: group=%s, stream=%s", logGroup, logStreamName)

	return ab.GetLogs(
		ctx,
		&params,
		func(gleo *cloudwatchlogs.GetLogEventsOutput, err error) error {
			if err != nil {
				return err
			}
			for _, event := range gleo.Events {
				log.Info(*event.Message)
			}
			return nil
		},
	)
}

// FIXME HACK HACK: we assume only one task and only one container in that task
func (ab *Backend) getLogEventsForTask(
	ctx context.Context,
	taskDefARN string,
	input *ecs.DescribeTasksInput,
	getlogs GetLogsFunc,
) error {
	err := ab.taskRunningWaiter.Wait(ctx, input, 600*time.Second)
	if err != nil {
		return errors.Wrap(err, "err waiting for tasks to start")
	}

	// get log groups
	taskDefResult, err := ab.ecsclient.DescribeTaskDefinition(
		ctx,
		&ecs.DescribeTaskDefinitionInput{TaskDefinition: &taskDefARN},
	)
	if err != nil {
		return errors.Wrap(err, "could not describe task definition")
	}

	// get log streams
	tasksResult, err := ab.ecsclient.DescribeTasks(ctx, input)
	if err != nil {
		return errors.Wrap(err, "could not describe tasks")
	}

	if tasksResult == nil || len(tasksResult.Tasks) == 0 || len(*tasksResult.Tasks[0].TaskArn) == 0 {
		return errors.Errorf("no matching tasks for task definition %s", taskDefARN)
	}
	taskId, err := ab.getTaskId(*tasksResult.Tasks[0].TaskArn)
	if err != nil {
		return errors.Wrap(err, "unable to determine a task id")
	}

	logGroup, logStreamName, err := ab.getLogGroupAndStreamName(*taskDefResult.TaskDefinition, taskId, "")
	if err != nil {
		return errors.Wrap(err, "unable to determine log group and stream name")
	}

	return ab.GetLogs(
		ctx,
		&cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &logGroup,
			LogStreamName: &logStreamName,
		},
		getlogs,
	)
}

func (ab *Backend) getTaskId(taskArn string) (string, error) {
	resourceArn, err := arn.Parse(taskArn)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse task ARN: '%s'", taskArn)
	}

	arnSegments := strings.Split(resourceArn.Resource, "/")
	if len(arnSegments) < 3 {
		return "", errors.Errorf("incomplete task ARN: '%s'", taskArn)
	}
	return arnSegments[len(arnSegments)-1], nil
}

func (ab *Backend) getLogGroupAndStreamName(taskDefinition ecstypes.TaskDefinition, taskId string, containerName string) (string, string, error) {
	logGroup := ""
	logStreamName := ""
	for _, containerDefinition := range taskDefinition.ContainerDefinitions {
		// If container name is specified, we only look at that container
		if len(containerName) > 0 && (*containerDefinition.Name != containerName) {
			continue
		}
		logGroup, ok := containerDefinition.LogConfiguration.Options[AwsLogsGroup]
		if !ok {
			continue
		}
		logStreamName := taskId
		logPrefix, ok := containerDefinition.LogConfiguration.Options[AwsLogsStreamPrefix]
		if ok {
			logStreamName = path.Join(logPrefix, *containerDefinition.Name, taskId)
		}
		return logGroup, logStreamName, nil
	}
	if len(logGroup) == 0 {
		return "", "", errors.Errorf("unable to determine a log group for task '%s'", taskId)
	}

	if len(logStreamName) == 0 {
		return "", "", errors.Errorf("unable to determine a log stream name for task '%s'", taskId)
	}
	return "", "", errors.Errorf("unable to determine a log stream name for container '%s'", containerName)
}
