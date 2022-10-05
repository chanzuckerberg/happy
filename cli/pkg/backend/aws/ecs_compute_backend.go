package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/diagnostics"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type ECSComputeBackend struct {
	Backend     *Backend
	HappyConfig *config.HappyConfig
}

func NewECSComputeBackend(ctx context.Context, happyConfig *config.HappyConfig, b *Backend) (interfaces.ComputeBackend, error) {
	return &ECSComputeBackend{
		Backend:     b,
		HappyConfig: happyConfig,
	}, nil
}

func (b *ECSComputeBackend) GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error) {
	secretId := b.HappyConfig.GetSecretId()
	out, err := b.Backend.secretsclient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretId,
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not get integration secret at %s", secretId)
	}

	secret := &config.IntegrationSecret{}
	err = json.Unmarshal([]byte(*out.SecretString), secret)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not json parse integraiton secret")
	}
	return secret, out.ARN, nil
}

func (b *ECSComputeBackend) GetParam(ctx context.Context, name string) (string, error) {
	logrus.Debugf("reading aws ssm parameter at %s", name)

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
func (b *ECSComputeBackend) GetLogGroupStreamsForStack(ctx context.Context, stackName string, serviceName string) (string, []string, error) {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "GetLogGroupStreamsForStack")

	stackTaskARNs, err := b.Backend.GetECSTasksForStackService(ctx, serviceName, stackName)
	if err != nil {
		return "", nil, errors.Wrapf(err, "error retrieving tasks for service '%s'", serviceName)
	}
	if len(stackTaskARNs) == 0 {
		return "", nil, errors.Errorf("no tasks associated with service %s. did you spell the service name correctly?", serviceName)
	}

	tasks, err := b.Backend.GetTaskDetails(ctx, stackTaskARNs)
	if err != nil {
		return "", nil, errors.Wrapf(err, "error getting task details for %+v", stackTaskARNs)
	}

	logConfigs, err := b.getAWSLogConfigsFromTasks(ctx, tasks...)
	if err != nil {
		return "", nil, err
	}

	return logConfigs.GroupName, logConfigs.StreamNames, nil
}

func (b *ECSComputeBackend) PrintLogs(ctx context.Context, stackName string, serviceName string, opts ...util.PrintOption) error {
	logGroup, logStreams, err := b.GetLogGroupStreamsForStack(ctx, stackName, serviceName)
	if err != nil {
		return err
	}
	p := util.MakeComputeLogPrinter(logGroup, logStreams, opts...)
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "PrintLogs")
	return p.Print(ctx, b.Backend.cwlFilterLogEventsAPIClient)
}

// RunTask runs an arbitrary task that is not necessarily associated with a service.
func (b *ECSComputeBackend) RunTask(ctx context.Context, taskDefArn string, launchType config.LaunchType) error {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), taskDefArn)

	log.Infof("running task %s", taskDefArn)
	out, err := b.Backend.ecsclient.RunTask(ctx, &ecs.RunTaskInput{
		Cluster:              &b.Backend.integrationSecret.ClusterArn,
		LaunchType:           ecstypes.LaunchType(launchType.String()),
		NetworkConfiguration: b.Backend.getNetworkConfig(),
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

	log.Infof("waiting for %+v to finish", tasks)
	err = b.Backend.waitForTasksToStop(ctx, tasks)
	if err != nil {
		return errors.Wrap(err, "error waiting for tasks")
	}

	log.Infof("task %s finished. printing logs from task", taskDefArn)
	logGroup, logStreams, err := b.getLogGroupStreamsForTasks(ctx, tasks)
	if err != nil {
		return err
	}

	p := util.MakeComputeLogPrinter(logGroup, logStreams, util.WithSince(util.GetStartTime(ctx).UnixMilli()))
	return p.Print(ctx, b.Backend.cwlFilterLogEventsAPIClient)
}

// GetLogGroupStreamsForTasks is just like GetLogGroupStreamsForService, except it gets all the log group and log
// streams associated with a one-off task, such as a migration or deletion task.
func (b *ECSComputeBackend) getLogGroupStreamsForTasks(ctx context.Context, taskARNs []string) (string, []string, error) {
	tasks, err := b.Backend.waitForTasksToStart(ctx, taskARNs)
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
			taskID, err := b.Backend.getTaskID(*task.TaskArn)
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
