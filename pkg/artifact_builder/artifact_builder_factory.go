package artifact_builder

import (
	"context"

	backend "github.com/chanzuckerberg/happy/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/profiler"
)

type ArtifactBuilderIface interface {
	WithConfig(config *BuilderConfig) ArtifactBuilderIface
	WithBackend(backend *backend.Backend) ArtifactBuilderIface
	WithTags(tags []string) ArtifactBuilderIface
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
	BuildAndPush(
		ctx context.Context,
		opts ...ArtifactBuilderBuildOption,
	) error
}

func NewArtifactBuilder(dryRun bool) ArtifactBuilderIface {
	if dryRun {
		return DryRunArtifactBuilder{}
	}
	return ArtifactBuilder{
		config:   nil,
		backend:  nil,
		Profiler: profiler.NewProfiler(),
		tags:     []string{},
	}
}
