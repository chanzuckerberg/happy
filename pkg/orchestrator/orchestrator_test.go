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
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"
const testDockerComposePath = "../config/testdata/docker-compose.yml"

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
	ecsApi.EXPECT().ListTasks(gomock.Any(), gomock.Any()).Return(&ecs.ListTasksOutput{}, nil)

	tasks := []types.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)
	containers := []types.Container{}
	containers = append(containers, types.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
		TaskArn:   aws.String("arn:::::ecs/task/name/mytaskid"),
	})
	tasks = append(tasks, types.Task{TaskArn: aws.String("arn:"),
		LastStatus:           aws.String("running"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           types.LaunchTypeEc2,
	})
	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{Tasks: tasks}, nil)

	containerInstances := []types.ContainerInstance{}
	containerInstances = append(containerInstances, types.ContainerInstance{Ec2InstanceId: aws.String("i-instance")})

	ecsApi.EXPECT().DescribeContainerInstances(gomock.Any(), gomock.Any()).Return(&ecs.DescribeContainerInstancesOutput{
		ContainerInstances: containerInstances,
	}, nil)
	ecsApi.EXPECT().RunTask(gomock.Any(), gomock.Any()).Return(&ecs.RunTaskOutput{
		Tasks: []types.Task{
			{LaunchType: types.LaunchTypeEc2},
		},
	}, nil)
	ecsApi.EXPECT().WaitUntilTasksRunning(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	ecsApi.EXPECT().WaitUntilTasksStopped(gomock.Any(), gomock.Any()).Return(nil)
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
							"awslogs-group":          "logsgroup",
							"awslogs--stream-prefix": "prefix-foobar",
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

	ecsApi.EXPECT().DescribeServices(gomock.Any(), gomock.Any()).Return(&ecs.DescribeServicesOutput{
		Services: []types.Service{
			{
				ServiceName: aws.String("name"),
				Deployments: []types.Deployment{
					{
						RolloutState: "PENDING",
					},
				},
				Events: []types.ServiceEvent{
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
	}, nil)

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
		}, nil)

	cwl := interfaces.NewMockGetLogEventsAPIClient(ctrl)
	cwl.EXPECT().GetLogEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cwlv2.GetLogEventsOutput{}, nil).Times(1)

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	req.NoError(err)

	backend, err := testbackend.NewBackend(
		ctx,
		ctrl,
		happyConfig,
		backend.WithECSClient(ecsApi),
		backend.WithEC2Client(ec2Api),
		backend.WithGetLogEventsAPIClient(cwl),
	)
	req.NoError(err)

	orchestrator := NewOrchestrator().WithBackend(backend).WithExecutor(util.NewDummyExecutor())
	req.NotNil(orchestrator)
	err = orchestrator.Shell(ctx, "frontend", "")
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

	stackMgr := stack_mgr.NewStackService().WithBackend(backend).WithWorkspaceRepo(mockWorkspaceRepo)
	stack := stack_mgr.NewStack(
		"stack1",
		stackMgr,
		util.NewLocalProcessor())
	err = orchestrator.RunTasks(ctx, stack, "delete")
	req.NoError(err)

	err = orchestrator.Logs("stack1", "frontend", time.Now().Add(time.Duration(-1)*time.Hour).String())
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
	ecsApi.EXPECT().ListTasks(gomock.Any(), gomock.Any()).Return(&ecs.ListTasksOutput{}, nil)

	tasks := []*types.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)
	containers := []types.Container{}
	containers = append(containers, types.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
	})
	tasks = append(tasks, &types.Task{TaskArn: aws.String("arn:"),
		LastStatus:           aws.String("running"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           types.LaunchTypeFargate,
	})
	ecsApi.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(&ecs.DescribeTasksOutput{Tasks: tasks}, nil)

	containerInstances := []types.ContainerInstance{}
	containerInstances = append(containerInstances, types.ContainerInstance{Ec2InstanceId: aws.String("i-instance")})

	ecsApi.EXPECT().DescribeContainerInstances(gomock.Any(), gomock.Any()).Return(&ecs.DescribeContainerInstancesOutput{
		ContainerInstances: containerInstances,
	}, nil)

	ecsApi.EXPECT().DescribeServices(gomock.Any(), gomock.Any()).Return(&ecs.DescribeServicesOutput{
		Services: []types.Service{
			{
				ServiceName: aws.String("name"),
				Deployments: []types.Deployment{
					{
						RolloutState: "PENDING",
					},
				},
				Events: []types.ServiceEvent{
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
	}, nil)

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
		}, nil)

	happyConfig, err := config.NewHappyConfig(bootstrapConfig)
	r.NoError(err)

	cwl := interfaces.NewMockGetLogEventsAPIClient(ctrl)

	backend, err := testbackend.NewBackend(
		ctx, ctrl, happyConfig,
		backend.WithECSClient(ecsApi),
		backend.WithEC2Client(ec2Api),
		backend.WithGetLogEventsAPIClient(cwl),
	)
	r.NoError(err)

	orchestrator := NewOrchestrator().WithBackend(backend).WithExecutor(util.NewDummyExecutor())
	r.NotNil(orchestrator)
	err = orchestrator.Shell(ctx, "frontend", "")
	r.NoError(err)

	err = orchestrator.GetEvents(ctx, "frontend", []string{"frontend"})
	r.NoError(err)
}
