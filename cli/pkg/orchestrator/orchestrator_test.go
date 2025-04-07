package orchestrator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cwlv2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwlv2types "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/chanzuckerberg/happy/cli/mocks"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/shared/config"
	stack_mgr "github.com/chanzuckerberg/happy/shared/stack"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/golang/mock/gomock"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../artifact_builder/testdata/test_config.yaml"
const testDockerComposePath = "../artifact_builder/testdata/docker-compose.yml"

func TestNewOrchestratorEC2(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.String())
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-RateLimit-Limit", "30")
		w.Header().Set("TFP-API-Version", "34.21.9")
		if r.URL.String() == "/api/v2/ping" {
			w.WriteHeader(204)
			return
		}

		fileName := fmt.Sprintf("./testdata%s.%s.json", r.URL.String(), r.Method)
		if strings.Contains(r.URL.String(), "/api/v2/state-version-outputs/") {
			fileName = fmt.Sprintf("./testdata%s.%s.json", "/api/v2/state-version-outputs", r.Method)
		}
		f, err := os.Open(fileName)
		req.NoError(err)
		_, err = io.Copy(w, f)
		req.NoError(err)

		w.WriteHeader(204)
	}))
	defer ts.Close()

	cf := &tfe.Config{
		Address:    ts.URL,
		Token:      "abcd1234",
		HTTPClient: ts.Client(),
	}

	client, err := tfe.NewClient(cf)
	if err != nil {
		t.Fatal(err)
	}

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "rdev",
	}

	ecsApi := interfaces.NewMockECSAPI(ctrl)

	ecsApi.EXPECT().ListServices(gomock.Any(), gomock.Any()).Return(&ecs.ListServicesOutput{
		ServiceArns: []string{
			"arn:aws:ecs:us-west-2:699936264352:task/blah/e7627b2daebe4744ab23fe36dba17739",
		},
	}, nil).AnyTimes()

	ecsApi.EXPECT().ListTasks(gomock.Any(), gomock.Any()).Return(&ecs.ListTasksOutput{
		NextToken: new(string),
		TaskArns:  []string{"arn:::::ecs/task/name/mytaskid"},
	}, nil).MaxTimes(5)

	tasks := []ecstypes.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)
	containers := []ecstypes.Container{}
	containers = append(containers, ecstypes.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
		TaskArn:   aws.String("arn:::::ecs/task/name/mytaskid"),
	})

	tasks = append(tasks, ecstypes.Task{TaskArn: aws.String("arn:::::ecs/task/name/mytaskid"),
		LastStatus:           aws.String("RUNNING"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           ecstypes.LaunchTypeEc2,
		TaskDefinitionArn:    aws.String("arn:aws:ecs:us-west-2:123456789012:task-definition/def:revision"),
	})
	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{Tasks: tasks}, nil).AnyTimes()

	taskStoppedWaiter := interfaces.NewMockECSTaskStoppedWaiterAPI(ctrl)
	taskStoppedWaiter.EXPECT().Wait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	containerInstances := []ecstypes.ContainerInstance{}
	containerInstances = append(containerInstances, ecstypes.ContainerInstance{Ec2InstanceId: aws.String("i-instance")})

	ecsApi.EXPECT().DescribeContainerInstances(gomock.Any(), gomock.Any()).Return(&ecs.DescribeContainerInstancesOutput{
		ContainerInstances: containerInstances,
	}, nil).AnyTimes()
	ecsApi.EXPECT().RunTask(gomock.Any(), gomock.Any()).Return(&ecs.RunTaskOutput{
		Tasks: []ecstypes.Task{
			{LaunchType: ecstypes.LaunchTypeEc2,
				TaskArn: aws.String("arn:::::ecs/task/name/mytaskid")},
		},
	}, nil)

	// ecsApi.EXPECT().WaitUntilTasksRunning(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	// ecsApi.EXPECT().WaitUntilTasksStopped(gomock.Any(), gomock.Any()).Return(nil)
	filterLogEventsApi := interfaces.NewMockFilterLogEventsAPIClient(ctrl)
	filterLogEventsApi.EXPECT().FilterLogEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cwlv2.FilterLogEventsOutput{}, nil).AnyTimes()

	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{
		Failures: []ecstypes.Failure{},
		Tasks:    tasks,
	}, nil).AnyTimes()

	ecsApi.EXPECT().DescribeTaskDefinition(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
		Tags: []ecstypes.Tag{},
		TaskDefinition: &ecstypes.TaskDefinition{
			TaskDefinitionArn: aws.String("arn:aws:ecs:us-west-2:123456789012:task-definition/def:revision"),
			Compatibilities:   []ecstypes.Compatibility{},
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
							backend.AwsLogsGroup:        "logsgroup",
							backend.AwsLogsStreamPrefix: "prefix-foobar",
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
			TaskRoleArn:             new(string),
			Volumes:                 []ecstypes.Volume{},
		},
	}, nil).AnyTimes()

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

	ec2Api := interfaces.NewMockEC2API(ctrl)

	ec2Api.EXPECT().DescribeInstances(gomock.Any(), gomock.Any()).Return(
		&ec2.DescribeInstancesOutput{Reservations: []ec2types.Reservation{
			{
				Groups: []ec2types.GroupIdentifier{},
				Instances: []ec2types.Instance{
					{
						PrivateIpAddress: aws.String("127.0.0.1"),
					},
				},
				OwnerId:       aws.String(""),
				RequesterId:   aws.String(""),
				ReservationId: aws.String(""),
			},
		},
		}, nil).AnyTimes()

	cwl := interfaces.NewMockGetLogEventsAPIClient(ctrl)
	cwl.EXPECT().GetLogEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cwlv2.GetLogEventsOutput{}, nil).AnyTimes()

	cwl.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).Return(&cwlv2.DescribeLogStreamsOutput{
		LogStreams: []cwlv2types.LogStream{
			{LogStreamName: aws.String("prefix-foobar/nginx/mytaskid")},
		},
		NextToken:      new(string),
		ResultMetadata: middleware.Metadata{},
	}, nil).AnyTimes()

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	req.NoError(err)

	backend, err := testbackend.NewBackend(
		ctx,
		ctrl,
		happyConfig.GetEnvironmentContext(),
		backend.WithECSClient(ecsApi),
		backend.WithEC2Client(ec2Api),
		backend.WithGetLogEventsAPIClient(cwl),
		backend.WithTaskStoppedWaiter(taskStoppedWaiter),
		backend.WithGetLogEventsAPIClient(cwl),
		backend.WithFilterLogEventsAPIClient(filterLogEventsApi),
		backend.WithExecutor(util.NewDummyExecutor()),
	)
	req.NoError(err)

	orchestrator := NewOrchestrator().WithHappyConfig(happyConfig).WithBackend(backend)
	req.NotNil(orchestrator)
	err = orchestrator.Shell(ctx, "frontend", "", "", []string{})
	req.NoError(err)

	err = orchestrator.GetEvents(ctx, "frontend", []string{"frontend"})
	req.NoError(err)

	mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
	ws := workspace_repo.TFEWorkspace{}
	ws.SetOutputs(map[string]string{"delete_db_task_definition_arn": "output"})
	currentRun := tfe.Run{ID: "run-CZcmD7eagjhyX0vN", ConfigurationVersion: &tfe.ConfigurationVersion{ID: "123"}}
	ws.SetClient(client)
	ws.SetWorkspace(&tfe.Workspace{ID: "workspace", CurrentRun: &currentRun})

	mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any(), gomock.Any()).Return(&ws, nil)

	stackMgr := stack_mgr.
		NewStackService(happyConfig.GetEnv(), happyConfig.App()).
		WithBackend(backend).
		WithWorkspaceRepo(mockWorkspaceRepo)
	stack := stack_mgr.NewStack(
		"stack1",
		stackMgr)

	req.True(orchestrator.TaskExists(ctx, "delete"))
	req.False(orchestrator.TaskExists(ctx, "create"))
	err = orchestrator.RunTasks(ctx, stack, "delete")
	req.NoError(err)

	cb, err := backend.GetComputeBackend(ctx)
	req.NoError(err)
	err = cb.PrintLogs(ctx, "stack1", "frontend", "")
	req.NoError(err)
}

func TestNewOrchestratorFargate(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	bootstrapConfig := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
		Env:                     "stage",
	}

	ecsApi := interfaces.NewMockECSAPI(ctrl)
	ecsApi.EXPECT().ListTasks(gomock.Any(), gomock.Any()).Return(&ecs.ListTasksOutput{
		NextToken: new(string),
		TaskArns:  []string{"arn:aws:ecs:us-east-1:123456789012:task/fargate-task-1"},
	}, nil)

	tasks := []ecstypes.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)
	containers := []ecstypes.Container{}
	containers = append(containers, ecstypes.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
	})

	tasks = append(tasks, ecstypes.Task{TaskArn: aws.String("arn:::::ecs/task/name/mytaskid"),
		LastStatus:           aws.String("running"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           ecstypes.LaunchTypeFargate,
	})
	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{Tasks: tasks}, nil).AnyTimes()

	containerInstances := []ecstypes.ContainerInstance{}
	containerInstances = append(containerInstances, ecstypes.ContainerInstance{Ec2InstanceId: aws.String("i-instance")})

	ecsApi.EXPECT().DescribeContainerInstances(gomock.Any(), gomock.Any()).Return(&ecs.DescribeContainerInstancesOutput{
		ContainerInstances: containerInstances,
	}, nil).AnyTimes()

	ecsApi.EXPECT().DescribeServices(gomock.Any(), gomock.Any()).Return(&ecs.DescribeServicesOutput{
		Services: []ecstypes.Service{
			{
				ServiceName: aws.String("name"),
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

	ec2Api := interfaces.NewMockEC2API(ctrl)
	ec2Api.EXPECT().DescribeInstances(gomock.Any(), gomock.Any()).Return(
		&ec2.DescribeInstancesOutput{Reservations: []ec2types.Reservation{
			{
				Groups: []ec2types.GroupIdentifier{},
				Instances: []ec2types.Instance{
					{
						PrivateIpAddress: aws.String("127.0.0.1"),
					},
				},
				OwnerId:       aws.String(""),
				RequesterId:   aws.String(""),
				ReservationId: aws.String(""),
			},
		},
		}, nil).AnyTimes()

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	cwl := interfaces.NewMockGetLogEventsAPIClient(ctrl)
	cwl.EXPECT().GetLogEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cwlv2.GetLogEventsOutput{}, nil).AnyTimes()

	cwl.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).Return(&cwlv2.DescribeLogStreamsOutput{
		LogStreams: []cwlv2types.LogStream{
			{LogStreamName: aws.String("123")},
		},
		NextToken:      new(string),
		ResultMetadata: middleware.Metadata{},
	}, nil).AnyTimes()

	backend, err := testbackend.NewBackend(
		ctx, ctrl, happyConfig.GetEnvironmentContext(),
		backend.WithECSClient(ecsApi),
		backend.WithEC2Client(ec2Api),
		backend.WithGetLogEventsAPIClient(cwl),
		backend.WithExecutor(util.NewDummyExecutor()),
	)
	r.NoError(err)

	orchestrator := NewOrchestrator().WithHappyConfig(happyConfig).WithBackend(backend)
	r.NotNil(orchestrator)
	err = orchestrator.Shell(ctx, "frontend", "", "", []string{})
	r.NoError(err)

	err = orchestrator.GetEvents(ctx, "frontend", []string{"frontend"})
	r.NoError(err)
}
