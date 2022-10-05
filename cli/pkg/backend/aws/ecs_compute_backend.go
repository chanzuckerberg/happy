package aws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	logConfigs, err := b.Backend.getAWSLogConfigsFromTasks(ctx, tasks...)
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
	logGroup, logStreams, err := b.Backend.GetLogGroupStreamsForTasks(ctx, tasks)
	if err != nil {
		return err
	}

	p := util.MakeComputeLogPrinter(logGroup, logStreams, util.WithSince(util.GetStartTime(ctx).UnixMilli()))
	return p.Print(ctx, b.Backend.cwlFilterLogEventsAPIClient)
}
