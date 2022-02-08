package config

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Environment struct {
	AWSProfile         string     `yaml:"aws_profile"`
	SecretARN          string     `yaml:"secret_arn"`
	TerraformDirectory string     `yaml:"terraform_directory"`
	DeleteProtected    bool       `yaml:"delete_protected"`
	AutoRunMigration   bool       `yaml:"auto_run_migration"`
	LogGroupPrefix     string     `yaml:"log_group_prefix"`
	TaskLaunchType     LaunchType `yaml:"task_launch_type"`
}

type ConfigData struct {
	ConfigVersion     string                 `yaml:"config_version"`
	TerraformVersion  string                 `yaml:"terraform_version"`
	DefaultEnv        string                 `yaml:"default_env"`
	App               string                 `yaml:"app"`
	DefaultComposeEnv string                 `yaml:"default_compose_env"`
	Environments      map[string]Environment `yaml:"environments"`
	Tasks             map[string][]string    `yaml:"tasks"`
	SliceDefaultTag   string                 `yaml:"slice_default_tag"`
	Slices            map[string]Slice       `yaml:"slices"`
	Services          []string               `yaml:"services"`
}

type Slice struct {
	BuildImages []string `yaml:"build_images"`
}

type HappyConfig struct {
	env  string
	data *ConfigData

	envConfig *Environment

	projectRoot string
	dockerRepo  string

	// serviceRegistries map[string]*RegistryConfig
	// clusterArn        string
	// privateSubnets    []string
	// securityGroups    []string
	// tfeUrl            string
	// tfeOrg            string
}

func NewHappyConfig(ctx context.Context, bootstrap *Bootstrap) (*HappyConfig, error) {
	configFilePath := bootstrap.GetHappyConfigPath()
	configContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}

	configData := &ConfigData{}
	err = yaml.Unmarshal(configContent, configData)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing yaml file")
	}

	env := bootstrap.GetEnv()
	envConfig, ok := configData.Environments[env]
	if !ok {
		return nil, errors.Errorf("environment not found: %s", env)
	}

	// awsProfile := envConfig.AWSProfile
	// if secretMgr == nil {
	// 	secretMgr = GetAwsSecretMgr(awsProfile)
	// }
	// secretArn := envConfig.SecretARN

	// secrets, err := secretMgr.GetSecrets(secretArn)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "unable to retrieve secrets")
	// }

	happyRootPath := bootstrap.GetHappyProjectRootPath()

	config := &HappyConfig{
		env:       env,
		data:      configData,
		envConfig: &envConfig,

		projectRoot: happyRootPath,

		// TODO: mv to integration secret
		// dockerRepo: dockerRepo,

		// serviceRegistries: secrets.GetAllServicesUrl(),
		// clusterArn:        secrets.GetClusterArn(),
		// privateSubnets:    secrets.GetPrivateSubnets(),
		// securityGroups:    secrets.GetSecurityGroups(),
		// tfeUrl:            secrets.GetTfeUrl(),
		// tfeOrg:            secrets.GetTfeOrg(),
	}

	return config, nil
}

func (s *HappyConfig) getData() *ConfigData {
	return s.data
}

func (s *HappyConfig) getEnvConfig() *Environment {
	return s.envConfig
}

func (s *HappyConfig) GetEnv() string {
	return s.env
}

func (s *HappyConfig) GetProjectRoot() string {
	return s.projectRoot
}

func (s *HappyConfig) AwsProfile() string {
	envConfig := s.getEnvConfig()

	return envConfig.AWSProfile
}

func (s *HappyConfig) GetSecretArn() string {
	envConfig := s.getEnvConfig()

	return envConfig.SecretARN
}

func (s *HappyConfig) AutoRunMigration() bool {
	envConfig := s.getEnvConfig()

	return envConfig.AutoRunMigration
}

func (s *HappyConfig) LogGroupPrefix() string {
	envConfig := s.getEnvConfig()

	return envConfig.LogGroupPrefix
}

func (s *HappyConfig) TerraformDirectory() string {
	envConfig := s.getEnvConfig()

	return envConfig.TerraformDirectory
}

func (s *HappyConfig) TaskLaunchType() LaunchType {
	envConfig := s.getEnvConfig()

	taskLaunchType := envConfig.TaskLaunchType
	if strings.ToUpper(taskLaunchType.String()) != LaunchTypeFargate.String() {
		taskLaunchType = LaunchTypeEC2
	}
	return LaunchTypeFargate
}

func (s *HappyConfig) TerraformVersion() string {
	return s.getData().TerraformVersion
}

func (s *HappyConfig) DefaultEnv() string {
	return s.getData().DefaultEnv
}

func (s *HappyConfig) DefaultComposeEnv() string {
	return s.getData().DefaultComposeEnv
}

func (s *HappyConfig) App() string {
	return s.getData().App
}

func (s *HappyConfig) GetTasks(taskType string) ([]string, error) {
	tasks, ok := s.getData().Tasks[taskType]
	if !ok {
		return nil, errors.Errorf("failed to get tasks: task type not found: %s", taskType)
	}
	return tasks, nil
}

func (s *HappyConfig) GetServices() []string {
	return s.getData().Services
}

// func (s *HappyConfig) GetRdevServiceRegistries() map[string]*RegistryConfig {
// 	return s.serviceRegistries
// }

// func (s *HappyConfig) ClusterArn() string {
// 	return s.clusterArn
// }

// func (s *HappyConfig) PrivateSubnets() []string {
// 	return s.privateSubnets
// }

// func (s *HappyConfig) SecurityGroups() []string {
// 	return s.securityGroups
// }

// func (s *HappyConfig) TfeUrl() string {
// 	return s.tfeUrl
// }

// func (s *HappyConfig) TfeOrg() string {
// 	return s.tfeOrg
// }

func (s *HappyConfig) SliceDefaultTag() string {
	return s.getData().SliceDefaultTag
}

func (s *HappyConfig) GetSlices() (map[string]Slice, error) {
	return s.getData().Slices, nil
}

func (s *HappyConfig) GetDockerRepo() string {
	return s.dockerRepo
}
