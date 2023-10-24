package artifact_builder

import (
	"context"

	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/profiler"
)

type ArtifactBuilderIface interface {
	WithConfig(config *BuilderConfig) ArtifactBuilderIface
	WithBackend(backend *backend.Backend) ArtifactBuilderIface
	WithHappyConfig(happyConfig *config.HappyConfig) ArtifactBuilderIface
	WithTags(tags []string) ArtifactBuilderIface
	GetTags() []string
	GetECRsForServices(ctx context.Context) (map[string]*config.RegistryConfig, error)
	CheckImageExists(ctx context.Context, tag string) (bool, error)
	RetagImages(
		ctx context.Context,
		serviceRegistries map[string]*config.RegistryConfig,
		sourceTag string,
		destTags []string,
		images []string,
	) error
	Build(ctx context.Context) error
	RegistryLogin(ctx context.Context) error
	Push(ctx context.Context, tags []string) error
	PushFromWithTag(ctx context.Context, servicesImage map[string]string, tag string) error
	Pull(ctx context.Context, stackName, tag string) (map[string]string, error)
	BuildAndPush(ctx context.Context) error
	DeleteImages(ctx context.Context, tag string) error
	GetServices(ctx context.Context) (map[string]ServiceConfig, error)
	GetAllServices(ctx context.Context) (map[string]ServiceConfig, error)
}

func CreateArtifactBuilder(ctx context.Context) ArtifactBuilderIface {
	return NewArtifactBuilder(ctx)
}

func NewArtifactBuilder(ctx context.Context) ArtifactBuilderIface {
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		return &DryRunArtifactBuilder{}
	}
	return &ArtifactBuilder{
		config:      nil,
		happyConfig: nil,
		backend:     nil,
		Profiler:    profiler.NewProfiler(),
		tags:        []string{},
	}
}
