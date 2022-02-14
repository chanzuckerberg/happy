package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/testbackend"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testFilePath = "../config/testdata/test_config.yaml"
const testDockerComposePath = "../config/testdata/docker-compose.yml"

func TestNewOrchestrator(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	ctrl := gomock.NewController(t)

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
	r.NoError(err)

	backend, err := testbackend.NewBackend(ctx, ctrl, happyConfig, backend.WithECSClient(ecsApi), backend.WithEC2Client(ec2Api))
	r.NoError(err)

	orchestrator := NewOrchestrator(backend)
	r.NotNil(orchestrator)
	err = orchestrator.Shell("frontend", "")
	r.NoError(err)

	err = orchestrator.GetEvents("frontend", []string{"frontend"})
	r.NoError(err)
}
