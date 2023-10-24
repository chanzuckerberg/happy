package artifact_builder

import (
	"context"
	"fmt"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ServiceBuild struct {
	Dockerfile string `yaml:"dockerfile"`
}

type ServiceConfig struct {
	Image    string                 `yaml:"image"`
	Build    *ServiceBuild          `yaml:"build"`
	Network  map[string]interface{} `yaml:"networks"`
	Platform string                 `yaml:"platform"`
}

type ConfigData struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

type BuilderConfig struct {
	composeFile    string
	composeEnvFile string
	dockerRepo     string
	env            string
	StackName      string
	Profile        *config.Profile

	// parse the passed in config file and populate some fields
	configData *ConfigData
	Executor   util.Executor
	DryRun     bool
}

func NewBuilderConfig() *BuilderConfig {
	return &BuilderConfig{
		Executor: util.NewDefaultExecutor(),
	}
}

func (b *BuilderConfig) Clone() *BuilderConfig {
	return &BuilderConfig{
		composeFile:    b.composeFile,
		composeEnvFile: b.composeEnvFile,
		dockerRepo:     b.dockerRepo,
		env:            b.env,
		StackName:      b.StackName,
		Profile:        b.Profile,
		configData:     b.configData,
		Executor:       b.Executor,
		DryRun:         b.DryRun,
	}
}

func (b *BuilderConfig) WithBootstrap(bootstrap *config.Bootstrap) *BuilderConfig {
	b.composeFile = bootstrap.DockerComposeConfigPath
	return b
}

func (b *BuilderConfig) WithHappyConfig(happyConfig *config.HappyConfig) *BuilderConfig {
	b.composeEnvFile = happyConfig.GetDockerComposeEnvFile()
	b.dockerRepo = happyConfig.GetDockerRepo()
	b.env = happyConfig.GetEnv()
	return b
}

func (s *BuilderConfig) GetContainers(ctx context.Context) ([]string, error) {
	var containers []string
	configData, err := s.retrieveConfigData(ctx)
	if err != nil {
		log.Errorf("unable to read config data: %s", err.Error())
		return containers, err
	}
	if configData.Services == nil {
		return containers, errors.New("no services defined in docker-compose.yml")
	}
	for _, service := range configData.Services {
		for _, network := range service.Network {
			for _, aliases := range network.(map[string]interface{}) {
				for _, alias := range aliases.([]interface{}) {
					containers = append(containers, alias.(string))
				}
			}
		}
	}

	return containers, nil
}

func (bc *BuilderConfig) retrieveConfigData(ctx context.Context) (*ConfigData, error) {
	if bc.configData != nil {
		return bc.configData, nil
	}

	configData, err := bc.DockerComposeConfig()
	if err != nil {
		return nil, err
	}
	bc.configData = configData
	err = bc.validateConfigData(ctx, configData)
	if err != nil {
		return nil, errors.Wrap(err, "unable to validate config data")
	}
	return bc.configData, nil
}

func (bc *BuilderConfig) validateConfigData(ctx context.Context, configData *ConfigData) error {
	for serviceName, service := range configData.Services {
		if len(service.Platform) == 0 {
			err := diagnostics.AddWarning(ctx, fmt.Sprintf("service '%s' has no platform defined in docker-compose.yaml which can lead to unexpected side effects", serviceName))
			if err != nil {
				return errors.Wrap(err, "unable to add warning")
			}
		}
	}
	return nil
}

func (s *BuilderConfig) GetConfigData(ctx context.Context) (*ConfigData, error) {
	if s.configData == nil {
		_, err := s.retrieveConfigData(ctx)
		if err != nil {
			return nil, err
		}
	}
	return s.configData, nil
}

// For testing purposes only
func (s *BuilderConfig) SetConfigData(configData *ConfigData) {
	s.configData = configData
}

func (s *BuilderConfig) GetBuildEnv() []string {
	dockerRepoStr := "DOCKER_REPO=" + s.dockerRepo

	return []string{
		fmt.Sprintf("HAPPY_ENV=%s", s.env),
		"DOCKER_BUILDKIT=1",
		"BUILDKIT_INLINE_CACHE=1",
		"COMPOSE_DOCKER_CLI_BUILD=1",
		dockerRepoStr,
	}
}

func (s *BuilderConfig) GetBuildServicesImage(ctx context.Context) (map[string]string, error) {
	configData, err := s.retrieveConfigData(ctx)
	if err != nil {
		return nil, err
	}

	svcs := map[string]string{}
	for serviceName, service := range configData.Services {
		// NOTE: we assume for now docker compose services without a build section are for local development only
		if service.Build == nil {
			log.Debugf("%s doesn't have a build section defined, skipping", serviceName)
			continue
		}
		svcs[serviceName] = service.Image
	}

	return svcs, nil
}
