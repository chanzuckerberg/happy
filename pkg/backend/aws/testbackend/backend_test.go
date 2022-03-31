package testbackend

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	awsbackend "github.com/chanzuckerberg/happy/pkg/backend/aws"
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

	tasks := []*ecs.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)
	containers := []*ecs.Container{}
	containers = append(containers, &ecs.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
		TaskArn:   aws.String("arn:::::ecs/task/name/mytaskid"),
	})
	tasks = append(tasks, &ecs.Task{
		LastStatus:           aws.String("running"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           aws.String("EC2"),
	})

	ecsApi := NewMockECSAPI(ctrl)
	ecsApi.EXPECT().RunTaskWithContext(gomock.Any(), gomock.Any()).Return(&ecs.RunTaskOutput{
		Tasks: []*ecs.Task{
			{LaunchType: aws.String("EC2")},
		},
	}, nil)
	ecsApi.EXPECT().WaitUntilTasksRunningWithContext(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	ecsApi.EXPECT().WaitUntilTasksStoppedWithContext(gomock.Any(), gomock.Any()).Return(nil)
	ecsApi.EXPECT().DescribeTasksWithContext(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{
		Failures: []*ecs.Failure{},
		Tasks:    tasks,
	}, nil).MaxTimes(5)
	ecsApi.EXPECT().DescribeTaskDefinitionWithContext(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
		Tags: []*ecs.Tag{},
		TaskDefinition: &ecs.TaskDefinition{
			Compatibilities: []*string{},
			ContainerDefinitions: []*ecs.ContainerDefinition{
				{
					Command:               []*string{},
					Cpu:                   new(int64),
					DependsOn:             []*ecs.ContainerDependency{},
					DisableNetworking:     new(bool),
					DnsSearchDomains:      []*string{},
					DnsServers:            []*string{},
					DockerLabels:          map[string]*string{},
					DockerSecurityOptions: []*string{},
					EntryPoint:            []*string{},
					Environment:           []*ecs.KeyValuePair{},
					EnvironmentFiles:      []*ecs.EnvironmentFile{},
					Essential:             new(bool),
					ExtraHosts:            []*ecs.HostEntry{},
					FirelensConfiguration: &ecs.FirelensConfiguration{},
					HealthCheck:           &ecs.HealthCheck{},
					Hostname:              new(string),
					Image:                 new(string),
					Interactive:           new(bool),
					Links:                 []*string{},
					LinuxParameters:       &ecs.LinuxParameters{},
					LogConfiguration: &ecs.LogConfiguration{
						Options: map[string]*string{
							"awslogs-group": aws.String("logsgroup"),
						},
					},
					Memory:                 new(int64),
					MemoryReservation:      new(int64),
					MountPoints:            []*ecs.MountPoint{},
					Name:                   new(string),
					PortMappings:           []*ecs.PortMapping{},
					Privileged:             new(bool),
					PseudoTerminal:         new(bool),
					ReadonlyRootFilesystem: new(bool),
					RepositoryCredentials:  &ecs.RepositoryCredentials{},
					ResourceRequirements:   []*ecs.ResourceRequirement{},
					Secrets:                []*ecs.Secret{},
					StartTimeout:           new(int64),
					StopTimeout:            new(int64),
					SystemControls:         []*ecs.SystemControl{},
					Ulimits:                []*ecs.Ulimit{},
					User:                   new(string),
					VolumesFrom:            []*ecs.VolumeFrom{},
					WorkingDirectory:       new(string),
				},
			},
			Cpu:                     new(string),
			DeregisteredAt:          &time.Time{},
			EphemeralStorage:        &ecs.EphemeralStorage{},
			ExecutionRoleArn:        new(string),
			Family:                  new(string),
			InferenceAccelerators:   []*ecs.InferenceAccelerator{},
			IpcMode:                 new(string),
			Memory:                  new(string),
			NetworkMode:             new(string),
			PidMode:                 new(string),
			PlacementConstraints:    []*ecs.TaskDefinitionPlacementConstraint{},
			ProxyConfiguration:      &ecs.ProxyConfiguration{},
			RegisteredAt:            &time.Time{},
			RegisteredBy:            new(string),
			RequiresAttributes:      []*ecs.Attribute{},
			RequiresCompatibilities: []*string{},
			Revision:                new(int64),
			RuntimePlatform:         &ecs.RuntimePlatform{},
			Status:                  new(string),
			TaskDefinitionArn:       new(string),
			TaskRoleArn:             new(string),
			Volumes:                 []*ecs.Volume{},
		},
	}, nil)

	cwl := NewMockGetLogEventsAPIClient(ctrl)
	cwl.EXPECT().GetLogEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cwlv2.GetLogEventsOutput{}, nil)

	b, err := NewBackend(ctx, ctrl, happyConfig, awsbackend.WithECSClient(ecsApi), awsbackend.WithGetLogEventsAPIClient(cwl))
	r.NoError(err)

	err = b.RunTask(ctx, "arn:task", "EC2")
	r.NoError(err)
}
