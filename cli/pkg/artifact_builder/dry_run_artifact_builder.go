package artifact_builder

import (
	"context"

	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DryRunArtifactBuilder struct {
	happyConfig *config.HappyConfig
	config      *BuilderConfig
}

// Build implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) Build(ctx context.Context) error {
	log.Info("Skipping Artifact Build")
	return nil
}

func (ab *DryRunArtifactBuilder) GetTags() []string {
	return []string{}
}

func (ab *DryRunArtifactBuilder) GetECRsForServices(ctx context.Context) (map[string]*config.RegistryConfig, error) {
	return nil, nil
}

func (ab *DryRunArtifactBuilder) Pull(ctx context.Context, stackName, tag string) (map[string]string, error) {
	return nil, nil
}

// BuildAndPush implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) BuildAndPush(ctx context.Context) error {
	log.Info("Skipping Artifact Build & Push")
	return nil
}

// CheckImageExists implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) CheckImageExists(ctx context.Context, tag string) (bool, error) {
	log.Info("Skipping Image Existence Check")
	return true, nil
}

// Push implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) Push(ctx context.Context, tags []string) error {
	log.Info("Skipping Artifact Push")
	return nil
}

// Push implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) PushFromWithTag(ctx context.Context, servicesImage map[string]string, tag string) error {
	log.Info("Skipping Artifact PushFrom")
	return nil
}

// RegistryLogin implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) RegistryLogin(ctx context.Context) error {
	log.Info("Skipping Registry Login")
	return nil
}

// RetagImages implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) RetagImages(ctx context.Context, serviceRegistries map[string]*config.RegistryConfig, sourceTag string, destTags []string, images []string) error {
	log.Info("Skipping Image Re-Tag")
	return nil
}

// WithBackend implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) WithBackend(backend *backend.Backend) ArtifactBuilderIface {
	return ab
}

// WithConfig implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) WithConfig(config *BuilderConfig) ArtifactBuilderIface {
	ab.config = config
	return ab
}

func (ab *DryRunArtifactBuilder) WithHappyConfig(happyConfig *config.HappyConfig) ArtifactBuilderIface {
	ab.happyConfig = happyConfig
	return ab
}

// WithTags implements ArtifactBuilderIface
func (ab *DryRunArtifactBuilder) WithTags(tags []string) ArtifactBuilderIface {
	return ab
}

func (ab *DryRunArtifactBuilder) GetServices(ctx context.Context) (map[string]ServiceConfig, error) {
	config, err := ab.config.GetConfigData(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get config data")
	}
	return config.Services, nil
}

func (ab *DryRunArtifactBuilder) GetAllServices(ctx context.Context) (map[string]ServiceConfig, error) {
	bc := BuilderConfig{
		composeFile:    ab.config.composeFile,
		composeEnvFile: ab.config.composeEnvFile,
		dockerRepo:     ab.config.dockerRepo,
		env:            ab.config.env,
		StackName:      ab.config.StackName,
		Profile:        nil,
		configData:     nil,
		Executor:       ab.config.Executor,
		DryRun:         ab.config.DryRun,
	}
	config, err := bc.GetConfigData(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get config data")
	}
	return config.Services, nil
}
