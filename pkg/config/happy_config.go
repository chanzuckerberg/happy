package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	LaunchTypeEC2     = "EC2"
	LaunchTypeFargate = "FARGATE"
)

type RegistryConfig struct {
	Url string `json:"url"`
}

func (s *RegistryConfig) GetRepoUrl() string {
	return s.Url
}

func (s *RegistryConfig) GetRegistryUrl() string {
	return strings.Split(s.Url, "/")[0]
}

type Environment struct {
	AWSProfile         string `yaml:"aws_profile"`
	SecretARN          string `yaml:"secret_arn"`
	TerraformDirectory string `yaml:"terraform_directory"`
	DeleteProtected    bool   `yaml:"delete_protected"`
	AutoRunMigration   bool   `yaml:"auto_run_migration"`
	LogGroupPrefix     string `yaml:"log_group_prefix"`
	TaskLaunchType     string `yaml:"task_launch_type"`
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

type HappyConfig interface {
	GetSecretArn() string
	GetProjectRoot() string
	GetTasks(taskType string) ([]string, error)
	AwsProfile() string
	AutoRunMigration() bool
	LogGroupPrefix() string
	TerraformDirectory() string
	TerraformVersion() string
	DefaultEnv() string
	DefaultComposeEnv() string
	App() string
	GetRdevServiceRegistries() (map[string]*RegistryConfig, error)
	ClusterArn() (string, error)
	PrivateSubnets() ([]string, error)
	SecurityGroups() ([]string, error)
	TfeUrl() (string, error)
	TfeOrg() (string, error)
	SliceDefaultTag() string
	GetSlices() (map[string]Slice, error)
	TaskLaunchType() string
	SetSecretsBackend(secretMgr SecretsBackend)
	GetServices() []string
	GetEnv() string
	GetDockerRepo() string
}

type happyConfig struct {
	env       string
	data      *ConfigData
	secretMgr SecretsBackend

	envConfig *Environment
	secrets   Secrets

	projectRoot string
	dockerRepo  string
}

func NewHappyConfig(bootstrap *Bootstrap) (HappyConfig, error) {
	configFilePath := bootstrap.GetHappyConfigPath()
	configContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}

	var configData ConfigData
	err = yaml.Unmarshal(configContent, &configData)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing yaml file")
	}

	env := bootstrap.GetEnv()
	envConfig, ok := configData.Environments[env]
	if !ok {
		return nil, errors.Errorf("environment not found: %s", env)
	}

	happyRootPath := bootstrap.GetHappyProjectRootPath()

	config := &happyConfig{
		env:       env,
		data:      &configData,
		envConfig: &envConfig,

		projectRoot: happyRootPath,
	}

	dockerRepo := os.Getenv("DOCKER_REPO")
	if len(dockerRepo) == 0 {
		serviceRegistries, err := config.GetRdevServiceRegistries()
		if err != nil {
			log.Errorf("Unable to retrieve registry information: %s\n", err.Error())
		} 
			for _, registry := range serviceRegistries {
				dockerRepo = registry.Url
				parts := strings.Split(registry.GetRepoUrl(), "/")
				if len(parts) < 2 {
					continue
				}
				dockerRepo = parts[0] + "/"
				break
			}
	}

	config.dockerRepo = dockerRepo

	return config, nil
}

func (s *happyConfig) getData() *ConfigData {
	return s.data
}

func (s *happyConfig) getEnvConfig() *Environment {
	return s.envConfig
}

func (s *happyConfig) GetEnv() string {
	return s.env
}

func (s *happyConfig) GetProjectRoot() string {
	return s.projectRoot
}

func (s *happyConfig) AwsProfile() string {
	envConfig := s.getEnvConfig()

	return envConfig.AWSProfile
}

func (s *happyConfig) GetSecretArn() string {
	envConfig := s.getEnvConfig()

	return envConfig.SecretARN
}

func (s *happyConfig) AutoRunMigration() bool {
	envConfig := s.getEnvConfig()

	return envConfig.AutoRunMigration
}

func (s *happyConfig) LogGroupPrefix() string {
	envConfig := s.getEnvConfig()

	return envConfig.LogGroupPrefix
}

func (s *happyConfig) TerraformDirectory() string {
	envConfig := s.getEnvConfig()

	return envConfig.TerraformDirectory
}

func (s *happyConfig) TaskLaunchType() string {
	envConfig := s.getEnvConfig()

	taskLaunchType := strings.ToUpper(envConfig.TaskLaunchType)
	if taskLaunchType != LaunchTypeFargate {
		taskLaunchType = LaunchTypeEC2
	}
	return taskLaunchType
}

func (s *happyConfig) TerraformVersion() string {
	return s.getData().TerraformVersion
}

func (s *happyConfig) DefaultEnv() string {
	return s.getData().DefaultEnv
}

func (s *happyConfig) DefaultComposeEnv() string {
	return s.getData().DefaultComposeEnv
}

func (s *happyConfig) App() string {
	return s.getData().App
}

func (s *happyConfig) GetTasks(taskType string) ([]string, error) {
	tasks, ok := s.getData().Tasks[taskType]
	if !ok {
		return nil, errors.Errorf("failed to get tasks: task type not found: %s", taskType)
	}
	return tasks, nil
}

func (s *happyConfig) GetServices() []string {
	return s.getData().Services
}

func (s *happyConfig) getSecrets() (Secrets, error) {
	if s.secretMgr == nil {
		awsProfile := s.AwsProfile()
		s.secretMgr = GetAwsSecretMgr(awsProfile)
	}

	secretArn := s.GetSecretArn()

	if s.secrets == nil {
		secrets, err := s.secretMgr.GetSecrets(secretArn)
		if err != nil {
			return nil, err
		}
		s.secrets = secrets
	}

	return s.secrets, nil
}

func (s *happyConfig) GetRdevServiceRegistries() (map[string]*RegistryConfig, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return nil, err
	}
	serviceRegistries := secrets.GetAllServicesUrl()
	return serviceRegistries, nil
}

func (s *happyConfig) ClusterArn() (string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return "", err
	}

	clusterArn := secrets.GetClusterArn()
	return clusterArn, nil
}

func (s *happyConfig) PrivateSubnets() ([]string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return nil, err
	}

	privateSubnets := secrets.GetPrivateSubnets()
	return privateSubnets, nil
}

func (s *happyConfig) SecurityGroups() ([]string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return nil, err
	}

	securityGroups := secrets.GetSecurityGroups()
	return securityGroups, nil
}

func (s *happyConfig) TfeUrl() (string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return "", err
	}

	tfeUrl := secrets.GetTfeUrl()
	return tfeUrl, nil
}

func (s *happyConfig) TfeOrg() (string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return "", err
	}

	tfeOrg := secrets.GetTfeOrg()
	return tfeOrg, nil
}

func (s *happyConfig) SliceDefaultTag() string {
	return s.getData().SliceDefaultTag
}

func (s *happyConfig) GetSlices() (map[string]Slice, error) {
	return s.getData().Slices, nil
}

// NOTE: testonly; TODO: add to linting rules to assert this
func (s *happyConfig) SetSecretsBackend(secretMgr SecretsBackend) {
	s.secretMgr = secretMgr
}

func (s *happyConfig) GetDockerRepo() string {

	return s.dockerRepo
}
