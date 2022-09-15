package interfaces

import (
	"context"

	"github.com/chanzuckerberg/happy/pkg/config"
)

type ComputeBackend interface {
	GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error)
}
