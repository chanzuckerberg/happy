package interfaces

import (
	"context"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
)

type StackServiceDescription struct {
	Compute string
	Params  map[string]string
}

type ComputeBackend interface {
	GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error)
	GetParam(ctx context.Context, name string) (string, error)
	WriteParam(ctx context.Context, name string, val string) error
	PrintLogs(ctx context.Context, stackName, serviceName, containerName string, opts ...util.PrintOption) error
	RunTask(ctx context.Context, taskDefArn string, launchType util.LaunchType) error
	Shell(ctx context.Context, stackName, serviceName, containerName string, shellCommand []string) error
	GetEvents(ctx context.Context, stackName string, services []string) error
	Describe(ctx context.Context, stackName string, serviceName string) (StackServiceDescription, error)
	GetResources(ctx context.Context, stackName string) ([]util.ManagedResource, error)
}
