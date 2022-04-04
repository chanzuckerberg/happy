package interfaces

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECSAPI interface {
	DescribeServices(ctx context.Context, params *ecs.DescribeServicesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error)
	ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
	DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	DescribeTaskDefinition(ctx context.Context, params *ecs.DescribeTaskDefinitionInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error)
	RunTask(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error)
	DescribeContainerInstances(ctx context.Context, params *ecs.DescribeContainerInstancesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeContainerInstancesOutput, error)
}

type ECSTaskRunningWaiterAPI interface {
	Wait(ctx context.Context, params *ecs.DescribeTasksInput, maxWaitDur time.Duration, optFns ...func(*ecs.TasksRunningWaiterOptions)) error
}

type ECSTaskStoppedWaiterAPI interface {
	Wait(ctx context.Context, params *ecs.DescribeTasksInput, maxWaitDur time.Duration, optFns ...func(*ecs.TasksStoppedWaiterOptions)) error
}
