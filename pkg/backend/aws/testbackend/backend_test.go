package testbackend

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	awsbackend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../../../config/testdata/test_config.yaml"
const testDockerComposePath = "../../../config/testdata/docker-compose.yml"

func TestAWSBackend(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

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
		TaskArn:              aws.String("arn:::::ecs/task/name/mytaskid"),
	})

	ecsApi := interfaces.NewMockECSAPI(ctrl)
	ecsApi.EXPECT().RunTask(gomock.Any(), gomock.Any()).Return(&ecs.RunTaskOutput{
		Tasks: []ecstypes.Task{
			{LaunchType: ecstypes.LaunchTypeEc2,
				TaskArn: aws.String("arn:::::ecs/task/name/mytaskid")},
		},
	}, nil)

	taskRunningWaiter := interfaces.NewMockECSTaskRunningWaiterAPI(ctrl)
	taskStoppedWaiter := interfaces.NewMockECSTaskStoppedWaiterAPI(ctrl)
	taskRunningWaiter.EXPECT().Wait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
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
					Name:                   new(string),
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
	cwl.EXPECT().GetLogEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cloudwatchlogs.GetLogEventsOutput{}, nil)

	b, err := NewBackend(ctx, ctrl, happyConfig,
		awsbackend.WithECSClient(ecsApi),
		awsbackend.WithGetLogEventsAPIClient(cwl),
		awsbackend.WithTaskRunningWaiter(taskRunningWaiter),
		awsbackend.WithTaskStoppedWaiter(taskStoppedWaiter))
	r.NoError(err)

	err = b.RunTask(ctx, "arn:::::ecs/task/name/mytaskid", "EC2")
	r.NoError(err)
}
