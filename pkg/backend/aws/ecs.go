package aws

import (
	"context"
	"path"
	"strings"

	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (b *Backend) RunTask(
	ctx context.Context,
	taskDefArn string,
	launchType config.LaunchType,
) error {
	clusterARN := b.integrationSecret.ClusterArn
	networkConfig := b.getNetworkConfig()

	out, err := b.ecsclient.RunTaskWithContext(ctx, &ecs.RunTaskInput{
		Cluster:              &clusterARN,
		LaunchType:           aws.String(launchType.String()),
		NetworkConfiguration: networkConfig,
		TaskDefinition:       &taskDefArn,
	})
	if err != nil {
		return errors.Wrapf(err, "could not run task %s", taskDefArn)
	}

	if len(out.Tasks) == 0 {
		return errors.New("could not run task, not found")
	}

	tasks := []*string{}
	for _, task := range out.Tasks {
		tasks = append(tasks, task.TaskArn)
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
			func(gleo *cwlv2.GetLogEventsOutput, err error) error {
				if err != nil {
					return err
				}
				select {
				case <-done:
					logrus.Warn("stopping")
					return Stop()
				default:
					for _, event := range gleo.Events {
						logrus.Info(*event.Message)
					}
					return nil
				}
			},
		)
	}()

	err = b.waitForTasks(ctx, waitInput)
	logrus.Warn("wait for tasks done")
	close(done)
	if err != nil {
		return errors.Wrap(err, "error waiting for tasks")
	}

	return <-logserr
}

func (ab *Backend) waitForTasks(ctx context.Context, input *ecs.DescribeTasksInput) error {
	// Wait until they are all running
	err := ab.ecsclient.WaitUntilTasksRunningWithContext(ctx, input)
	if err != nil {
		return errors.Wrap(err, "err waiting for tasks to start")
	}

	// Wait until they all stop
	err = ab.ecsclient.WaitUntilTasksStoppedWithContext(ctx, input)
	if err != nil {
		return errors.Wrap(err, "err waiting for tasks to stop")
	}

	// now get their status
	tasks, err := ab.ecsclient.DescribeTasksWithContext(ctx, input)
	if err != nil {
		return errors.Wrap(err, "could not describe tasks")
	}

	var failures error
	for _, failure := range tasks.Failures {
		failures = multierror.Append(failures, errors.Errorf("error running task (%s) with status (%s) and reason (%s)", *failure.Arn, *failure.Detail, *failure.Reason))
	}
	return failures
}

func (ab *Backend) getNetworkConfig() *ecs.NetworkConfiguration {
	privateSubnets := ab.integrationSecret.PrivateSubnets
	privateSubnetsPt := []*string{}
	for _, subnet := range privateSubnets {
		subnetValue := subnet
		privateSubnetsPt = append(privateSubnetsPt, &subnetValue)
	}
	securityGroups := ab.integrationSecret.SecurityGroups
	securityGroupsPt := []*string{}
	for _, sg := range securityGroups {
		sgValue := sg
		securityGroupsPt = append(securityGroupsPt, &sgValue)
	}

	awsvpcConfiguration := &ecs.AwsVpcConfiguration{
		AssignPublicIp: aws.String("DISABLED"),
		SecurityGroups: securityGroupsPt,
		Subnets:        privateSubnetsPt,
	}
	networkConfig := &ecs.NetworkConfiguration{
		AwsvpcConfiguration: awsvpcConfiguration,
	}
	return networkConfig
}

// FIXME HACK HACK: we assume only one task and only one container in that task
func (ab *Backend) getLogEventsForTask(
	ctx context.Context,
	taskDefARN string,
	input *ecs.DescribeTasksInput,
	getlogs GetLogsFunc,
) error {
	// Wait until they are all running
	err := ab.ecsclient.WaitUntilTasksRunningWithContext(ctx, input)
	if err != nil {
		return errors.Wrap(err, "err waiting for tasks to start")
	}

	// get log groups
	taskDefResult, err := ab.ecsclient.DescribeTaskDefinitionWithContext(
		ctx,
		&ecs.DescribeTaskDefinitionInput{TaskDefinition: &taskDefARN},
	)
	if err != nil {
		return errors.Wrap(err, "could not describe task definition")
	}

	// get log streams
	tasksResult, err := ab.ecsclient.DescribeTasksWithContext(ctx, input)
	if err != nil {
		return errors.Wrap(err, "could not describe tasks")
	}

	// NOTE NOTE: we are making an assumption that we only have one container per task
	//            this was here before, but I don't know if it is valid
	container := tasksResult.Tasks[0].Containers[0]
	if container.Reason != nil {
		logrus.Warnf("container exited with status %s: %s", *container.LastStatus, *container.Reason)
	}

	// now the log group
	logConfiguration := taskDefResult.TaskDefinition.ContainerDefinitions[0].LogConfiguration
	logGroup, ok := logConfiguration.Options["awslogs-group"]
	if !ok {
		return errors.Errorf("could not infer log group")
	}

	logPrefix, logPrefixSpecified := logConfiguration.Options["awslogs-stream-prefix"]
	// Now we should have enough info to sort out the log stream
	// see https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html
	// NOTE: this is REQUIRED for fargate, but shares logic with ECS as well
	//       when prefix defined
	taskARN := strings.Split(*container.TaskArn, "/")
	logStream := taskARN[len(taskARN)-1]
	if logPrefixSpecified {
		// prefix-name/container-name/ecs-task-id
		prefixName := *logPrefix
		containerName := *taskDefResult.TaskDefinition.ContainerDefinitions[0].Name
		ecsTaskID := logStream
		logStream = path.Join(prefixName, containerName, ecsTaskID)
	}

	return ab.getLogs(
		ctx,
		&cwlv2.GetLogEventsInput{
			LogGroupName:  logGroup,
			LogStreamName: &logStream,
		},
		getlogs,
	)
}
