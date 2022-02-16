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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/chanzuckerberg/happy/mocks"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
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

	ecsApi := testbackend.NewMockECSAPI(ctrl)
	ecsApi.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{}, nil)

	tasks := []*ecs.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)
	containers := []*ecs.Container{}
	containers = append(containers, &ecs.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
	})
	tasks = append(tasks, &ecs.Task{TaskArn: aws.String("arn:"),
		LastStatus:           aws.String("running"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           aws.String("EC2"),
	})
	ecsApi.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{Tasks: tasks}, nil)

	containerInstances := []*ecs.ContainerInstance{}
	containerInstances = append(containerInstances, &ecs.ContainerInstance{Ec2InstanceId: aws.String("i-instance")})

	ecsApi.EXPECT().DescribeContainerInstances(gomock.Any()).Return(&ecs.DescribeContainerInstancesOutput{
		ContainerInstances: containerInstances,
	}, nil)

	ecsApi.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
		Services: []*ecs.Service{
			{
				ServiceName: aws.String("name"),
				Deployments: []*ecs.Deployment{
					{
						RolloutState: aws.String("PENDING"),
					},
				},
				Events: []*ecs.ServiceEvent{
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

	ec2Api := testbackend.NewMockEC2API(ctrl)
	ec2Api.EXPECT().DescribeInstances(gomock.Any()).Return(
		&ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{
			{
				Groups: []*ec2.GroupIdentifier{},
				Instances: []*ec2.Instance{
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

	happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
	req.NoError(err)

	backend, err := testbackend.NewBackend(ctx, ctrl, happyConfig, backend.WithECSClient(ecsApi), backend.WithEC2Client(ec2Api))
	req.NoError(err)

	orchestrator := NewOrchestrator(backend).WithExecutor(util.NewDummyExecutor())
	req.NotNil(orchestrator)
	err = orchestrator.Shell("frontend", "")
	req.NoError(err)

	err = orchestrator.GetEvents("frontend", []string{"frontend"})
	req.NoError(err)

	mockWorkspaceRepo := mocks.NewMockWorkspaceRepoIface(ctrl)
	ws := workspace_repo.TFEWorkspace{}
	currentRun := tfe.Run{ID: "run-CZcmD7eagjhyX0vN", ConfigurationVersion: &tfe.ConfigurationVersion{ID: "123"}}
	ws.SetClient(client)
	ws.SetWorkspace(&tfe.Workspace{ID: "workspace", CurrentRun: &currentRun})

	mockWorkspaceRepo.EXPECT().GetWorkspace(gomock.Any()).Return(&ws, nil)

	stackMgr := stack_mgr.NewStackService(backend, mockWorkspaceRepo)
	stack := stack_mgr.NewStack(
		"stack1",
		stackMgr,
		util.NewLocalProcessor())
	err = orchestrator.RunTasks(ctx, stack, "delete", false)
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

	ecsApi := testbackend.NewMockECSAPI(ctrl)
	ecsApi.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{}, nil)

	tasks := []*ecs.Task{}
	startedAt := time.Now().Add(time.Duration(-2) * time.Hour)
	containers := []*ecs.Container{}
	containers = append(containers, &ecs.Container{
		Name:      aws.String("nginx"),
		RuntimeId: aws.String("123"),
	})
	tasks = append(tasks, &ecs.Task{TaskArn: aws.String("arn:"),
		LastStatus:           aws.String("running"),
		ContainerInstanceArn: aws.String("host"),
		StartedAt:            &startedAt,
		Containers:           containers,
		LaunchType:           aws.String("FARGATE"),
	})
	ecsApi.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{Tasks: tasks}, nil)

	containerInstances := []*ecs.ContainerInstance{}
	containerInstances = append(containerInstances, &ecs.ContainerInstance{Ec2InstanceId: aws.String("i-instance")})

	ecsApi.EXPECT().DescribeContainerInstances(gomock.Any()).Return(&ecs.DescribeContainerInstancesOutput{
		ContainerInstances: containerInstances,
	}, nil)

	ecsApi.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
		Services: []*ecs.Service{
			{
				ServiceName: aws.String("name"),
				Deployments: []*ecs.Deployment{
					{
						RolloutState: aws.String("PENDING"),
					},
				},
				Events: []*ecs.ServiceEvent{
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

	ec2Api := testbackend.NewMockEC2API(ctrl)
	ec2Api.EXPECT().DescribeInstances(gomock.Any()).Return(
		&ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{
			{
				Groups: []*ec2.GroupIdentifier{},
				Instances: []*ec2.Instance{
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

	happyConfig, err := config.NewHappyConfig(ctx, bootstrapConfig)
	r.NoError(err)

	backend, err := testbackend.NewBackend(ctx, ctrl, happyConfig, backend.WithECSClient(ecsApi), backend.WithEC2Client(ec2Api))
	r.NoError(err)

	orchestrator := NewOrchestrator(backend).WithExecutor(util.NewDummyExecutor())
	r.NotNil(orchestrator)
	err = orchestrator.Shell("frontend", "")
	r.NoError(err)

	err = orchestrator.GetEvents("frontend", []string{"frontend"})
	r.NoError(err)
}
