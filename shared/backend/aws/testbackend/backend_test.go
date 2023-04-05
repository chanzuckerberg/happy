package testbackend

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	awsbackend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../../../config/testdata/test_config.yaml"
const testDockerComposePath = "../../../config/testdata/docker-compose.yml"

func TestAWSBackend(t *testing.T) {
	r := require.New(t)

	ctx := context.WithValue(context.Background(), util.CmdStartContextKey, time.Now())

	ctrl := gomock.NewController(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	tasks := []ecstypes.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)

	containers := []ecstypes.Container{}
	containers = append(containers, ecstypes.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
		TaskArn:   aws.String("arn:::::ecs/task/name/mytaskid"),
	})

	tasks = append(tasks, ecstypes.Task{
		LastStatus:           aws.String("RUNNING"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           ecstypes.LaunchTypeEc2,
		TaskDefinitionArn:    aws.String("arn:aws:ecs:us-west-2:123456789012:task-definition/hello_world:8"),
		TaskArn:              aws.String("arn:::::ecs/task/name/mytaskid"),
	})

	ecsApi := interfaces.NewMockECSAPI(ctrl)
	ecsApi.EXPECT().RunTask(gomock.Any(), gomock.Any()).Return(&ecs.RunTaskOutput{
		Tasks: []ecstypes.Task{
			{LaunchType: ecstypes.LaunchTypeEc2,
				TaskDefinitionArn: aws.String("arn:aws:ecs:us-west-2:123456789012:task-definition/hello_world:8"),
				TaskArn:           aws.String("arn:::::ecs/task/name/mytaskid")},
		},
	}, nil)

	cloudwatchApi := interfaces.NewMockGetLogEventsAPIClient(ctrl)
	cloudwatchApi.EXPECT().GetLogEvents(gomock.Any(), gomock.Any()).Return(&cloudwatchlogs.GetLogEventsOutput{}, nil).AnyTimes()
	cloudwatchApi.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []cloudwatchtypes.LogStream{
			{LogStreamName: aws.String("123")},
		},
		NextToken:      new(string),
		ResultMetadata: middleware.Metadata{},
	}, nil).AnyTimes()

	filterLogEventsApi := interfaces.NewMockFilterLogEventsAPIClient(ctrl)
	filterLogEventsApi.EXPECT().FilterLogEvents(gomock.Any(), gomock.Any()).Return(&cloudwatchlogs.FilterLogEventsOutput{}, nil).AnyTimes()

	taskStoppedWaiter := interfaces.NewMockECSTaskStoppedWaiterAPI(ctrl)
	taskStoppedWaiter.EXPECT().Wait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{
		Failures: []ecstypes.Failure{},
		Tasks:    tasks,
	}, nil).MaxTimes(5)
	ecsApi.EXPECT().DescribeTaskDefinition(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
		Tags: []ecstypes.Tag{},
		TaskDefinition: &ecstypes.TaskDefinition{
			Compatibilities: []ecstypes.Compatibility{},
			ContainerDefinitions: []ecstypes.ContainerDefinition{
				{
					Command:               []string{},
					Cpu:                   0,
					DependsOn:             []ecstypes.ContainerDependency{},
					DisableNetworking:     new(bool),
					DnsSearchDomains:      []string{},
					DnsServers:            []string{},
					DockerLabels:          map[string]string{},
					DockerSecurityOptions: []string{},
					EntryPoint:            []string{},
					Environment:           []ecstypes.KeyValuePair{},
					EnvironmentFiles:      []ecstypes.EnvironmentFile{},
					Essential:             new(bool),
					ExtraHosts:            []ecstypes.HostEntry{},
					FirelensConfiguration: &ecstypes.FirelensConfiguration{},
					HealthCheck:           &ecstypes.HealthCheck{},
					Hostname:              new(string),
					Image:                 new(string),
					Interactive:           new(bool),
					Links:                 []string{},
					LinuxParameters:       &ecstypes.LinuxParameters{},
					LogConfiguration: &ecstypes.LogConfiguration{
						Options: map[string]string{
							awsbackend.AwsLogsGroup: "logsgroup",
						},
					},
					Memory:                 new(int32),
					MemoryReservation:      new(int32),
					MountPoints:            []ecstypes.MountPoint{},
					Name:                   aws.String("nginx"),
					PortMappings:           []ecstypes.PortMapping{},
					Privileged:             new(bool),
					PseudoTerminal:         new(bool),
					ReadonlyRootFilesystem: new(bool),
					RepositoryCredentials:  &ecstypes.RepositoryCredentials{},
					ResourceRequirements:   []ecstypes.ResourceRequirement{},
					Secrets:                []ecstypes.Secret{},
					StartTimeout:           new(int32),
					StopTimeout:            new(int32),
					SystemControls:         []ecstypes.SystemControl{},
					Ulimits:                []ecstypes.Ulimit{},
					User:                   new(string),
					VolumesFrom:            []ecstypes.VolumeFrom{},
					WorkingDirectory:       new(string),
				},
			},
			Cpu:                     new(string),
			DeregisteredAt:          &time.Time{},
			EphemeralStorage:        &ecstypes.EphemeralStorage{},
			ExecutionRoleArn:        new(string),
			Family:                  new(string),
			InferenceAccelerators:   []ecstypes.InferenceAccelerator{},
			IpcMode:                 "",
			Memory:                  new(string),
			NetworkMode:             "",
			PidMode:                 "",
			PlacementConstraints:    []ecstypes.TaskDefinitionPlacementConstraint{},
			ProxyConfiguration:      &ecstypes.ProxyConfiguration{},
			RegisteredAt:            &time.Time{},
			RegisteredBy:            new(string),
			RequiresAttributes:      []ecstypes.Attribute{},
			RequiresCompatibilities: []ecstypes.Compatibility{},
			Revision:                0,
			RuntimePlatform:         &ecstypes.RuntimePlatform{},
			Status:                  "",
			TaskDefinitionArn:       new(string),
			TaskRoleArn:             new(string),
			Volumes:                 []ecstypes.Volume{},
		},
	}, nil)

	cwl := interfaces.NewMockGetLogEventsAPIClient(ctrl)
	cwl.EXPECT().GetLogEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cloudwatchlogs.GetLogEventsOutput{}, nil).AnyTimes()

	b, err := NewBackend(ctx, ctrl, happyConfig.GetEnvironmentContext(),
		awsbackend.WithECSClient(ecsApi),
		awsbackend.WithGetLogEventsAPIClient(cwl),
		awsbackend.WithTaskStoppedWaiter(taskStoppedWaiter),
		awsbackend.WithGetLogEventsAPIClient(cloudwatchApi),
		awsbackend.WithFilterLogEventsAPIClient(filterLogEventsApi),
	)
	r.NoError(err)

	err = b.RunTask(ctx, "arn:::::ecs/task/name/mytaskid", "EC2")
	r.NoError(err)
}
