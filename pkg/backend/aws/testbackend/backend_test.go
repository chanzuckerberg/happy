package testbackend

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
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

	tasks := []types.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)

	containers := []types.Container{}
	containers = append(containers, types.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
		TaskArn:   aws.String("arn:::::ecs/task/name/mytaskid"),
	})

	tasks = append(tasks, types.Task{
		LastStatus:           aws.String("RUNNING"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           types.LaunchTypeEc2,
		TaskArn:              aws.String("arn:::::ecs/task/name/mytaskid"),
	})

	ecsApi := interfaces.NewMockECSAPI(ctrl)
	ecsApi.EXPECT().RunTask(gomock.Any(), gomock.Any()).Return(&ecs.RunTaskOutput{
		Tasks: []types.Task{
			{LaunchType: types.LaunchTypeEc2,
				TaskArn: aws.String("arn:::::ecs/task/name/mytaskid")},
		},
	}, nil)

	taskRunningWaiter := interfaces.NewMockECSTaskRunningWaiterAPI(ctrl)
	taskStoppedWaiter := interfaces.NewMockECSTaskStoppedWaiterAPI(ctrl)
	taskRunningWaiter.EXPECT().Wait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	taskStoppedWaiter.EXPECT().Wait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{
		Failures: []types.Failure{},
		Tasks:    tasks,
	}, nil).MaxTimes(5)
	ecsApi.EXPECT().DescribeTaskDefinition(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
		Tags: []types.Tag{},
		TaskDefinition: &types.TaskDefinition{
			Compatibilities: []types.Compatibility{},
			ContainerDefinitions: []types.ContainerDefinition{
				{
					Command:               []string{},
					Cpu:                   0,
					DependsOn:             []types.ContainerDependency{},
					DisableNetworking:     new(bool),
					DnsSearchDomains:      []string{},
					DnsServers:            []string{},
					DockerLabels:          map[string]string{},
					DockerSecurityOptions: []string{},
					EntryPoint:            []string{},
					Environment:           []types.KeyValuePair{},
					EnvironmentFiles:      []types.EnvironmentFile{},
					Essential:             new(bool),
					ExtraHosts:            []types.HostEntry{},
					FirelensConfiguration: &types.FirelensConfiguration{},
					HealthCheck:           &types.HealthCheck{},
					Hostname:              new(string),
					Image:                 new(string),
					Interactive:           new(bool),
					Links:                 []string{},
					LinuxParameters:       &types.LinuxParameters{},
					LogConfiguration: &types.LogConfiguration{
						Options: map[string]string{
							awsbackend.AwsLogsGroup: "logsgroup",
						},
					},
					Memory:                 new(int32),
					MemoryReservation:      new(int32),
					MountPoints:            []types.MountPoint{},
					Name:                   new(string),
					PortMappings:           []types.PortMapping{},
					Privileged:             new(bool),
					PseudoTerminal:         new(bool),
					ReadonlyRootFilesystem: new(bool),
					RepositoryCredentials:  &types.RepositoryCredentials{},
					ResourceRequirements:   []types.ResourceRequirement{},
					Secrets:                []types.Secret{},
					StartTimeout:           new(int32),
					StopTimeout:            new(int32),
					SystemControls:         []types.SystemControl{},
					Ulimits:                []types.Ulimit{},
					User:                   new(string),
					VolumesFrom:            []types.VolumeFrom{},
					WorkingDirectory:       new(string),
				},
			},
			Cpu:                     new(string),
			DeregisteredAt:          &time.Time{},
			EphemeralStorage:        &types.EphemeralStorage{},
			ExecutionRoleArn:        new(string),
			Family:                  new(string),
			InferenceAccelerators:   []types.InferenceAccelerator{},
			IpcMode:                 "",
			Memory:                  new(string),
			NetworkMode:             "",
			PidMode:                 "",
			PlacementConstraints:    []types.TaskDefinitionPlacementConstraint{},
			ProxyConfiguration:      &types.ProxyConfiguration{},
			RegisteredAt:            &time.Time{},
			RegisteredBy:            new(string),
			RequiresAttributes:      []types.Attribute{},
			RequiresCompatibilities: []types.Compatibility{},
			Revision:                0,
			RuntimePlatform:         &types.RuntimePlatform{},
			Status:                  "",
			TaskDefinitionArn:       new(string),
			TaskRoleArn:             new(string),
			Volumes:                 []types.Volume{},
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
