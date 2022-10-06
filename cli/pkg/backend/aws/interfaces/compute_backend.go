package interfaces

import (
	"context"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
)

type ComputeBackend interface {
	GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error)
	GetParam(ctx context.Context, name string) (string, error)
	WriteParam(ctx context.Context, name string, val string) error
	PrintLogs(ctx context.Context, stackName string, serviceName string, opts ...util.PrintOption) error
	RunTask(ctx context.Context, taskDefArn string, launchType config.LaunchType) error
}
