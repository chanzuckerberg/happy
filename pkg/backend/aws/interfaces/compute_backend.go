package interfaces

import (
	"context"

	"github.com/chanzuckerberg/happy/pkg/config"
)

type ComputeBackend interface {
	GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error)
	GetParam(ctx context.Context, name string) (string, error)
	WriteParam(ctx context.Context, name string, val string) error
}
