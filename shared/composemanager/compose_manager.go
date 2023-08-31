package composemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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

type ComposeProject struct {
	Version  string                   `yaml:"version,omitempty" json:"version,omitempty"`
	Services map[string]ServiceConfig `yaml:"services,omitempty" json:"services,omitempty"`
}

type ServiceConfig struct {
	Image    string            `yaml:"image,omitempty" json:"image,omitempty"`
	Platform string            `yaml:"platform,omitempty" json:"platform,omitempty"`
	Profiles []string          `yaml:"profiles,omitempty" json:"profiles,omitempty"`
	Build    types.BuildConfig `yaml:"build,omitempty" json:"build,omitempty"`
	Ports    []string          `yaml:"ports,omitempty" json:"ports,omitempty"`
}

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

func (c ComposeManager) Manage(ctx context.Context) error {
	composeFilePath := c.HappyConfig.GetBootstrap().DockerComposeConfigPath
	if _, err := os.Stat(composeFilePath); err == nil || !errors.Is(err, os.ErrNotExist) {
		var choice string
		prompt := []*survey.Question{
			{
				Name: "environments",
				Prompt: &survey.Select{
					Message: fmt.Sprintf("File %s already exists. Would you like to overwrite it, save a backup, or skip the change?", composeFilePath),
					Options: []string{"overwrite", "backup", "skip"},
					Default: "skip",
				},
			},
		}

		err := survey.Ask(prompt, &choice, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrapf(err, "failed to prompt")
		}

		if choice == "skip" {
			return nil
		} else if choice == "overwrite" {
			return c.Save(ctx)
		} else if choice == "backup" {
			err = os.Rename(composeFilePath, fmt.Sprintf("%s.bak", composeFilePath))
			if err != nil {
				return errors.Wrapf(err, "unable to move %s to %s.bak", composeFilePath, composeFilePath)
			}
			logrus.Infof("You can find the existing file at %s.bak, you might need to manually merge the files.", composeFilePath)
			return c.Save(ctx)
		} else {
			return errors.New("invalid choice")
		}
	}

	return c.Save(ctx)
}

func (c ComposeManager) Save(ctx context.Context) error {
	composeFilePath := c.HappyConfig.GetBootstrap().DockerComposeConfigPath

	p := ComposeProject{
		Version: "3.8",
	}
	p.Services = map[string]ServiceConfig{}

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
		sd, ok := servicesDef[service]
		if !ok {
			continue
		}

		var servicePort uint32
		var buildConfig map[string]any
		switch m := sd.(type) {
		case map[string]any:
			port, ok := m["port"].(float64)
			if !ok {
				logrus.Warnf("service definition for '%s' does not have a port specified, skipping", service)
				continue
			}
			servicePort = uint32(port)
			if err != nil {
				logrus.Warnf("service definition for '%s' does not have a valid port specified, skipping", service)
				continue
			}
			buildConfig, ok = m["build"].(map[string]any)
			if !ok {
				logrus.Warnf("service definition for '%s' does not have a build config specificed, skipping", service)
				continue
			}
		default:
			logrus.Warnf("service definition for '%s' is not a string map", service)
			continue
		}

		serviceDef := sd.(map[string]any)
		serviceConfig := ServiceConfig{
			Image:    service,
			Profiles: []string{"*"},
			Build:    types.BuildConfig{},
			Ports:    []string{fmt.Sprintf("%d:%d", servicePort, servicePort)},
		}

		platform, ok := serviceDef[platform_architecture].(string)
		if !ok || len(platform) == 0 {
			platform = "arm64"
		}

		jsonData, err := json.Marshal(buildConfig)
		if err != nil {
			return errors.Wrap(err, "unable to marshal build config")
		}

		err = json.Unmarshal(jsonData, &serviceConfig.Build)
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal build config")
		}

		serviceConfig.Platform = fmt.Sprintf("linux/%s", platform)
		p.Services[service] = serviceConfig
	}

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

	p, err := loadComposeFile(c.HappyConfig.GetBootstrap().HappyProjectRoot, composeFilePath)
	if err != nil {
		return errors.Wrapf(err, "unable to load compose file %s", composeFilePath)
	}

	stackDef, err := c.HappyConfig.GetStackConfig()
	if err != nil {
		return errors.Wrap(err, "unable to get stack config")
	}

	_, ok := stackDef[services]
	if !ok {
		return errors.New("unable to find services in stack config")
	}

	switch stackDef[services].(type) {
	case map[string]any:
	default:
		return errors.New("invalid happy config file structure: services are not configured as a string map")
	}

	servicesDef, ok := stackDef[services].(map[string]any)
	if !ok || len(servicesDef) == 0 {
		return errors.New("no service settings are defined in stack config")
	}

	composeServiceMap := map[string]types.ServiceConfig{}
	for _, service := range p.Services {
		composeServiceMap[service.Name] = service
		if _, ok := servicesDef[service.Name]; !ok {
			return errors.Errorf("service '%s' from docker-compose is not defined in stack config", service.Name)
		}
	}

	for serviceName := range servicesDef {
		var composeServiceDef types.ServiceConfig
		var ok bool

		if composeServiceDef, ok = composeServiceMap[serviceName]; !ok {
			continue
		}
		serviceDef := servicesDef[serviceName].(map[string]any)
		serviceDef[build] = composeServiceDef.Build

		composePlatformArchitecture := ""
		if len(composeServiceDef.Platform) > 0 {
			composePlatformArchitecture = composeServiceDef.Platform
		}

		platformArchitecture := ""
		if arch, ok := serviceDef[platform_architecture]; ok {
			if len(platformArchitecture) > 0 {
				platformArchitecture = arch.(string)
			}
		}

		if len(composePlatformArchitecture) > 0 && len(platformArchitecture) > 0 && composePlatformArchitecture != fmt.Sprintf("linux/%s", platformArchitecture) {
			return errors.Errorf("platform_architecture mismatch for service %s", serviceName)
		}

		if len(composePlatformArchitecture) > 0 {
			serviceDef[platform_architecture] = normalizeArchitecture(composePlatformArchitecture)
		} else if len(platformArchitecture) > 0 {
			serviceDef[platform_architecture] = normalizeArchitecture(platformArchitecture)
		} else {
			serviceDef[platform_architecture] = "arm64"
		}
	}
	c.HappyConfig.GetData().StackDefaults[services] = servicesDef
	return errors.Wrap(c.HappyConfig.Save(), "unable to save happy config")
}

func loadComposeFile(projectRoot, composeFilePath string) (*types.Project, error) {
	configYaml, err := os.ReadFile(composeFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read %s", composeFilePath)
	}

	p, err := loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename: composeFilePath,
				Content:  configYaml,
			},
		},
		Environment: map[string]string{},
		WorkingDir:  projectRoot,
	}, func(o *loader.Options) {
		o.SetProjectName("happy", false)
		o.SkipNormalization = true
	}, loader.WithProfiles([]string{"*"}))

	if err != nil {
		return nil, errors.Wrap(err, "unable to load compose config")
	}
	return p, nil
}

func normalizeArchitecture(arch string) string {
	arch = strings.TrimPrefix(arch, "linux/")
	if arch == "amd64" || arch == "arm64" {
		return arch
	}
	return "arm64"
}
