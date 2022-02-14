package artifact_builder

import (
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ServiceBuild struct {
	Dockerfile string `yaml:"dockerfile"`
}

type ServiceConfig struct {
	Image   string                 `yaml:"image"`
	Build   *ServiceBuild          `yaml:"build"`
	Network map[string]interface{} `yaml:"networks"`
}

type ConfigData struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

type BuilderConfig struct {
	composeFile    string
	composeEnvFile string
	dockerRepo     string

	profile *config.Profile

	// parse the passed in config file and populate some fields
	configData *ConfigData
}

func NewBuilderConfig(bootstrap *config.Bootstrap, happyConfig *config.HappyConfig) *BuilderConfig {
	return &BuilderConfig{
		composeFile:    bootstrap.DockerComposeConfigPath,
		composeEnvFile: happyConfig.GetDockerComposeEnvFile(),
		dockerRepo:     happyConfig.GetDockerRepo(),
	}
}

func (b *BuilderConfig) WithProfile(p *config.Profile) *BuilderConfig {
	b.profile = p
	return b
}

func (s *BuilderConfig) GetContainers() ([]string, error) {
	var containers []string
	configData, err := s.getConfigData()
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

func (bc *BuilderConfig) getConfigData() (*ConfigData, error) {
	if bc.configData != nil {
		return bc.configData, nil
	}

	configData, err := bc.DockerComposeConfig()
	if err != nil {
		return nil, err
	}
	bc.configData = configData
	return bc.configData, nil
}

func (s *BuilderConfig) GetBuildEnv() []string {
	dockerRepoStr := "DOCKER_REPO=" + s.dockerRepo

	return []string{
		"DOCKER_BUILDKIT=1",
		"BUILDKIT_INLINE_CACHE=1",
		"COMPOSE_DOCKER_CLI_BUILD=1",
		dockerRepoStr,
	}
}

func (s *BuilderConfig) GetBuildServicesImage() (map[string]string, error) {
	configData, err := s.getConfigData()
	if err != nil {
		return nil, err
	}

	svcs := map[string]string{}
	for serviceName, service := range configData.Services {
		if service.Build != nil {
			svcs[serviceName] = service.Image
		}
	}
	return svcs, nil
}
