package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go/middleware"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	compute "github.com/chanzuckerberg/happy/shared/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNetworkConfig(t *testing.T) {
	r := require.New(t)
	backend := Backend{}
	sgs := []string{"sg-1", "sg-2"}
	subnets := []string{"subnet-1", "subnet-2"}

	backend.integrationSecret = &config.IntegrationSecret{
		ClusterArn:     "arn:cluster",
		PrivateSubnets: subnets,
		SecurityGroups: sgs,
		Services:       map[string]*config.RegistryConfig{},
	}
	ecsBackend := ECSComputeBackend{Backend: &backend}
	networkConfig := ecsBackend.getNetworkConfig()
	r.NotNil(networkConfig)
	r.Equal(len(subnets), len(networkConfig.AwsvpcConfiguration.Subnets))
	r.Equal(len(sgs), len(networkConfig.AwsvpcConfiguration.SecurityGroups))

	for index, subnet := range subnets {
		r.Equal(subnet, networkConfig.AwsvpcConfiguration.Subnets[index])
	}
	for index, sg := range sgs {
		r.Equal(sg, networkConfig.AwsvpcConfiguration.SecurityGroups[index])
	}
}

func TestEcsTasks(t *testing.T) {
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

	secretsApi := interfaces.NewMockSecretsManagerAPI(ctrl)
	testVal := "{\"cluster_arn\": \"test_arn\",\"ecrs\": {\"ecr_1\": {\"url\": \"test_url_1\"}},\"tfe\": {\"url\": \"tfe_url\",\"org\": \"tfe_org\"}}"
	secretsApi.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: &testVal,
			ARN:          aws.String("arn:aws:secretsmanager:region:accountid:secret:happy/env-happy-config-AB1234"),
		}, nil).AnyTimes()

	stsApi := interfaces.NewMockSTSAPI(ctrl)
	stsApi.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).
		Return(&sts.GetCallerIdentityOutput{UserId: aws.String("foo:bar")}, nil).AnyTimes()

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
	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{
		Failures: []ecstypes.Failure{},
		Tasks:    tasks,
	}, nil).AnyTimes()
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
						LogDriver: "awslogs",
						Options: map[string]string{
							AwsLogsGroup: "logsgroup",
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
	}, nil).AnyTimes()
	ecsApi.EXPECT().RunTask(gomock.Any(), gomock.Any()).Return(&ecs.RunTaskOutput{
		Tasks: []ecstypes.Task{
			{LaunchType: ecstypes.LaunchTypeEc2,
				TaskDefinitionArn: aws.String("arn:aws:ecs:us-west-2:123456789012:task-definition/hello_world:8"),
				TaskArn:           aws.String("arn:::::ecs/task/name/mytaskid")},
		},
	}, nil).AnyTimes()
	ecsApi.EXPECT().ListServices(gomock.Any(), gomock.Any()).Return(&ecs.ListServicesOutput{
		ServiceArns: []string{
			"arn:aws:ecs:us-west-2:123456789012:task/blah/e7627b2daebe4744ab23fe36dba17739",
		},
	}, nil).AnyTimes()

	ecsApi.EXPECT().ListTasks(gomock.Any(), gomock.Any()).Return(&ecs.ListTasksOutput{
		NextToken: new(string),
		TaskArns:  []string{"arn:aws:ecs:us-west-2:123456789012:task/fargate-task-1"},
	}, nil)
	ecsApi.EXPECT().DescribeServices(gomock.Any(), gomock.Any()).Return(&ecs.DescribeServicesOutput{
		Services: []ecstypes.Service{
			{
				ServiceName: aws.String("stack1-frontend"),
				Deployments: []ecstypes.Deployment{
					{
						RolloutState: "PENDING",
					},
				},
				Events: []ecstypes.ServiceEvent{
					{
						CreatedAt: &startedAt,
						Message:   aws.String("deregistered"),
					},
					{
						CreatedAt: &startedAt,
						Message:   aws.String("deregistered"),
					},
					{
						CreatedAt: &startedAt,
						Message:   aws.String("deregistered"),
					},
					{
						CreatedAt: &startedAt,
						Message:   aws.String("deregistered"),
					},
				},
			},
		},
	}, nil).AnyTimes()

	taskStoppedWaiter := interfaces.NewMockECSTaskStoppedWaiterAPI(ctrl)
	taskStoppedWaiter.EXPECT().Wait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

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

	computeBackend := compute.NewMockComputeBackend(ctrl)

	b, err := NewAWSBackend(ctx, happyConfig.GetEnvironmentContext(),
		WithAWSAccountID("123456789012"),
		WithSTSClient(stsApi),
		WithSecretsClient(secretsApi),
		WithECSClient(ecsApi),
		WithTaskStoppedWaiter(taskStoppedWaiter),
		WithGetLogEventsAPIClient(cloudwatchApi),
		WithFilterLogEventsAPIClient(filterLogEventsApi),
		WithComputeBackend(computeBackend),
	)
	r.NoError(err)

	ecsBackend := ECSComputeBackend{Backend: b}
	_, err = ecsBackend.GetTaskDefinitions(ctx, []string{"arn:::::ecs/task/name/mytaskid"})
	r.NoError(err)
	err = b.RunTask(ctx, "arn:::::ecs/task/name/mytaskid", "EC2")
	r.NoError(err)
	err = b.ComputeBackend.PrintLogs(ctx, "stack1", "frontend")
	r.NoError(err)
	taskId, err := ecsBackend.getTaskID("arn:::::ecs/task/name/mytaskid")
	r.NoError(err)
	r.Equal("mytaskid", taskId)
}
