package composemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	platform_architecture = "platform_architecture"
	build                 = "build"
	services              = "services"
)

type ComposeManager struct {
	HappyConfig *config.HappyConfig
}

func NewComposeManager() ComposeManager {
	return ComposeManager{}
}

func (c ComposeManager) WithHappyConfig(happyConfig *config.HappyConfig) ComposeManager {
	c.HappyConfig = happyConfig
	return c
}

func (c ComposeManager) Generate(ctx context.Context) error {
	p := types.Project{}
	p.Services = types.Services{}

	stackDef, err := c.HappyConfig.GetStackConfig()
	if err != nil {
		return errors.Wrap(err, "unable to get stack config")
	}

	_, ok := stackDef[services]
	if !ok {
		return errors.New("unable to find services in stack config")
	}

	servicesDef := stackDef[services].(map[string]any)
	if len(servicesDef) == 0 {
		return errors.New("no service settings are defined in stack config")
	}

	for _, service := range c.HappyConfig.GetData().Services {
		if sd, ok := servicesDef[service]; ok {
			serviceDef := sd.(map[string]any)
			serviceConfig := types.ServiceConfig{
				Name:     service,
				Image:    service,
				Profiles: []string{"*"},
				Build:    &types.BuildConfig{},
			}
			platform := serviceDef[platform_architecture].(string)
			if len(platform) == 0 {
				platform = "amd64"
			}

			jsonData, err := json.Marshal(serviceDef[build])
			if err != nil {
				return errors.Wrap(err, "unable to marshal build config")
			}

			err = json.Unmarshal(jsonData, serviceConfig.Build)
			if err != nil {
				return errors.Wrap(err, "unable to unmarshal build config")
			}

			serviceConfig.Platform = fmt.Sprintf("linux/%s", platform)
			p.Services = append(p.Services, serviceConfig)
		}
	}

	composeFilePath := c.HappyConfig.GetBootstrap().DockerComposeConfigPath
	logrus.Debugf("Generating docker-compose.yml at %s", composeFilePath)
	configYaml, err := yaml.Marshal(p)
	if err != nil {
		return errors.Wrap(err, "unable to marshal compose config")
	}

	err = os.WriteFile(composeFilePath, configYaml, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "unable to write out %s", composeFilePath)
	}
	return nil
}

func (c ComposeManager) Ingest(ctx context.Context) error {
	composeFilePath := c.HappyConfig.GetBootstrap().DockerComposeConfigPath
	logrus.Debugf("Ingesting docker-compose.yml at %s", composeFilePath)

	configYaml, err := os.ReadFile(composeFilePath)
	if err != nil {
		return errors.Wrapf(err, "unable to read %s", composeFilePath)
	}

	p, err := loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename: composeFilePath,
				Content:  configYaml,
			},
		},
		Environment: map[string]string{},
		WorkingDir:  c.HappyConfig.GetBootstrap().HappyProjectRoot,
	}, func(o *loader.Options) {
		o.SetProjectName("happy", false)
		o.SkipNormalization = true
	}, loader.WithProfiles([]string{"*"}))

	if err != nil {
		return errors.Wrap(err, "unable to load compose config")
	}

	stackDef, err := c.HappyConfig.GetStackConfig()
	if err != nil {
		return errors.Wrap(err, "unable to get stack config")
	}

	_, ok := stackDef[services]
	if !ok {
		return errors.New("unable to find services in stack config")
	}

	servicesDef := stackDef[services].(map[string]any)
	if len(servicesDef) == 0 {
		return errors.New("no service settings are defined in stack config")
	}

	composeServiceMap := map[string]types.ServiceConfig{}
	for _, service := range p.Services {
		composeServiceMap[service.Name] = service
	}

	for serviceName := range servicesDef {
		if composeServiceDef, ok := composeServiceMap[serviceName]; ok {
			serviceDef := servicesDef[serviceName].(map[string]any)
			serviceDef[build] = composeServiceDef.Build

			composePlatformArchitecture := "linux/amd64"
			if len(composeServiceDef.Platform) > 0 {
				composePlatformArchitecture = composeServiceDef.Platform
			}

			platformArchitecture := "amd64"
			if arch, ok := serviceDef[platform_architecture]; ok {
				if len(platformArchitecture) > 0 {
					platformArchitecture = arch.(string)
				}
			}

			if composePlatformArchitecture != fmt.Sprintf("linux/%s", platformArchitecture) {
				return errors.Errorf("platform_architecture mismatch for service %s", serviceName)
			}
			serviceDef[platform_architecture] = platformArchitecture
		}
	}
	c.HappyConfig.GetData().StackDefaults[services] = servicesDef
	return errors.Wrap(c.HappyConfig.Save(), "unable to save happy config")
}
