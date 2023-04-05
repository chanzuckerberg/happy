package artifact_builder

import (
	"context"

	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/profiler"
)

type ArtifactBuilderIface interface {
	WithConfig(config *BuilderConfig) ArtifactBuilderIface
	WithBackend(backend *backend.Backend) ArtifactBuilderIface
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
	BuildAndPush(ctx context.Context) error
}

func CreateArtifactBuilder(happyConfig *config.HappyConfig) ArtifactBuilderIface {
	return NewArtifactBuilder(happyConfig, false)
}

func NewArtifactBuilder(happyConfig *config.HappyConfig, dryRun bool) ArtifactBuilderIface {
	if bool(dryRun) {
		return DryRunArtifactBuilder{}
	}
	return &ArtifactBuilder{
		config:      nil,
		happyConfig: happyConfig,
		backend:     nil,
		Profiler:    profiler.NewProfiler(),
		tags:        []string{},
	}
}
