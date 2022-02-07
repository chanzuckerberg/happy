package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (ab *awsBackend) RunTask(
	ctx context.Context,
	taskDefArn string,
	launchType LaunchType,
) error {
	clusterARN := ab.conf.ClusterArn()
	networkConfig := ab.getNetworkConfig()

	out, err := ab.ecsclient.RunTaskWithContext(ctx, &ecs.RunTaskInput{
		Cluster:              &clusterARN,
		LaunchType:           aws.String(string(launchType)),
		NetworkConfiguration: networkConfig,
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

	// at this point, we want the log messages no matter what
	defer func() {
		// messages, err := ab.getLogs(ctx, )
	}()

	err = ab.waitForTasks(ctx, waitInput)
	if err != nil {
		return errors.Wrap(err, "error waiting for tasks")
	}

	// TODO
	return nil
}

func (ab *awsBackend) waitForTasks(ctx context.Context, input *ecs.DescribeTasksInput) error {
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

func (ab *awsBackend) getNetworkConfig() *ecs.NetworkConfiguration {
	privateSubnets := ab.conf.PrivateSubnets()
	privateSubnetsPt := []*string{}
	for _, subnet := range privateSubnets {
		privateSubnetsPt = append(privateSubnetsPt, &subnet)
	}
	securityGroups := ab.conf.SecurityGroups()
	securityGroupsPt := []*string{}
	for _, subnet := range securityGroups {
		securityGroupsPt = append(securityGroupsPt, &subnet)
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

func (ab *awsBackend) getLogEvents(
	ctx context.Context,
	taskDefARN string,
	launchType LaunchType,
	describeTasksInput *ecs.DescribeTasksInput,
) (*LogMessages, error) {
	// get log groups
	taskDefResult, err := ab.ecsclient.DescribeTaskDefinitionWithContext(
		ctx,
		&ecs.DescribeTaskDefinitionInput{TaskDefinition: &taskDefARN},
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not describe task definition")
	}

	// get log streams
	tasksResult, err := ab.ecsclient.DescribeTasksWithContext(ctx, describeTasksInput)
	if err != nil {
		return nil, errors.Wrap(err, "could not describe tasks")
	}

	// NOTE NOTE: we are making an assumption that we only have one container per task
	//            this was here before but I don't know if it is valid
	container := tasksResult.Tasks[0].Containers[0]
	if container.Reason != nil {
		logrus.Warnf("container exited with status %s: %s", *container.LastStatus, *container.Reason)
	}

	// now the log group
	logConfiguration := taskDefResult.TaskDefinition.ContainerDefinitions[0].LogConfiguration
	logGroup, ok := logConfiguration.Options["awslogs-group"]
	if !ok {
		return nil, errors.Errorf("could infer log group")
	}
	logPrefix := logConfiguration.Options["awslogs-stream-prefix"]

	logStream := container.RuntimeId
	if launchType == LaunchTypeFargate {

	}
	// TODO: sort this out

}
