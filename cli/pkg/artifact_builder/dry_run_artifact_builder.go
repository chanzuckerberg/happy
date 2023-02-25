package artifact_builder

import (
	"context"

	backend "github.com/chanzuckerberg/happy/cli/pkg/backend/aws"
	"github.com/chanzuckerberg/happy/cli/pkg/config"
	log "github.com/sirupsen/logrus"
)

type DryRunArtifactBuilder struct {
}

// Build implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) Build(ctx context.Context) error {
	log.Info("Skipping Artifact Build")
	return nil
}

func (ab DryRunArtifactBuilder) GetTags() []string {
	return []string{}
}

// BuildAndPush implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) BuildAndPush(ctx context.Context, opts ...ArtifactBuilderBuildOption) error {
	log.Info("Skipping Artifact Build & Push")
	return nil
}

// CheckImageExists implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) CheckImageExists(ctx context.Context, tag string) (bool, error) {
	log.Info("Skipping Image Existence Check")
	return true, nil
}

// Push implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) Push(ctx context.Context, tags []string) error {
	log.Info("Skipping Artifact Push")
	return nil
}

// RegistryLogin implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) RegistryLogin(ctx context.Context) error {
	log.Info("Skipping Registry Login")
	return nil
}

// RetagImages implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) RetagImages(ctx context.Context, serviceRegistries map[string]*config.RegistryConfig, sourceTag string, destTags []string, images []string) error {
	log.Info("Skipping Image Re-Tag")
	return nil
}

// WithBackend implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) WithBackend(backend *backend.Backend) ArtifactBuilderIface {
	return ab
}

// WithConfig implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) WithConfig(config *BuilderConfig) ArtifactBuilderIface {
	return ab
}

// WithTags implements ArtifactBuilderIface
func (ab DryRunArtifactBuilder) WithTags(tags []string) ArtifactBuilderIface {
	return ab
}
