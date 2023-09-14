package aws

import (
	"context"

	"github.com/chanzuckerberg/happy/shared/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
)

type NullComputeBackend struct {
	Backend *Backend
}

func NewNullComputeBackend(ctx context.Context, b *Backend) (*NullComputeBackend, error) {
	return &NullComputeBackend{
		Backend: b,
	}, nil
}

func (b *NullComputeBackend) GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error) {
	return &config.IntegrationSecret{
		Tfe: &config.TfeSecret{
			Org: "",
			Url: "",
		},
	}, nil, nil
}

func (b *NullComputeBackend) GetParam(ctx context.Context, name string) (string, error) {
	return "", errors.New("not implemented")
}

func (b *NullComputeBackend) WriteParam(
	ctx context.Context,
	name string,
	val string,
) error {
	return errors.New("not implemented")
}

func (b *NullComputeBackend) PrintLogs(ctx context.Context, stackName, serviceName, containerName string, opts ...util.PrintOption) error {
	return errors.New("not implemented")
}

func (b *NullComputeBackend) RunTask(ctx context.Context, taskDefArn string, launchType util.LaunchType) error {
	return errors.New("not implemented")
}

func (b *NullComputeBackend) Shell(ctx context.Context, stackName, service, containerName, shellCommand string) error {
	return errors.New("not implemented")
}

func (b *NullComputeBackend) GetEvents(ctx context.Context, stackName string, services []string) error {
	return errors.New("not implemented")
}

func (b *NullComputeBackend) Describe(ctx context.Context, stackName string, serviceName string) (interfaces.StackServiceDescription, error) {
	return interfaces.StackServiceDescription{}, errors.New("not implemented")
}

func (b *NullComputeBackend) GetResources(ctx context.Context, stackName string) ([]util.ManagedResource, error) {
	return []util.ManagedResource{}, errors.New("not implemented")
}

func (b *NullComputeBackend) ListClusterIds(ctx context.Context) ([]string, error) {
	return []string{}, errors.New("not implemented")
}
