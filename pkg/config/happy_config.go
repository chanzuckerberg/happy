package config

import (
	"io/ioutil"
	"strings"

	// artifactBuilder "github.com/chanzuckerberg/happy/pkg/artifact_builder"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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
}

type Slice struct {
	BuildImages []string `yaml:"build_images"`
}

type HappyConfigIface interface {
	GetSecretArn() string
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
}

type HappyConfig struct {
	env       string
	data      *ConfigData
	secretMgr SecretsBackend

	envConfig *Environment
	secrets   Secrets
}

func NewHappyConfig(configFile string, env string) (HappyConfigIface, error) {
	configFile, err := homedir.Expand(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse aws config file path")
	}

	configContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read file")
	}

	var configData ConfigData
	err = yaml.Unmarshal(configContent, &configData)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing yaml file")
	}

	envConfig, ok := configData.Environments[env]
	if !ok {
		return nil, errors.Errorf("Environment not found: %s", env)
	}

	return &HappyConfig{
		env:       env,
		data:      &configData,
		envConfig: &envConfig,
	}, err
}

func (s *HappyConfig) getData() *ConfigData {
	return s.data
}

func (s *HappyConfig) getEnvConfig() *Environment {
	return s.envConfig
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

func (s *HappyConfig) getSecrets() (Secrets, error) {

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

func (s *HappyConfig) GetRdevServiceRegistries() (map[string]*RegistryConfig, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return nil, err
	}
	serviceRegistries := secrets.GetAllServicesUrl()
	return serviceRegistries, nil
}

func (s *HappyConfig) ClusterArn() (string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return "", err
	}

	clusterArn := secrets.GetClusterArn()
	return clusterArn, nil
}

func (s *HappyConfig) PrivateSubnets() ([]string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return nil, err
	}

	privateSubnets := secrets.GetPrivateSubnets()
	return privateSubnets, nil
}

func (s *HappyConfig) SecurityGroups() ([]string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return nil, err
	}

	securityGroups := secrets.GetSecurityGroups()
	return securityGroups, nil
}

func (s *HappyConfig) TfeUrl() (string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return "", err
	}

	tfeUrl := secrets.GetTfeUrl()
	return tfeUrl, nil
}

func (s *HappyConfig) TfeOrg() (string, error) {
	secrets, err := s.getSecrets()
	if err != nil {
		return "", err
	}

	tfeOrg := secrets.GetTfeOrg()
	return tfeOrg, nil
}

func (s *HappyConfig) SliceDefaultTag() string {
	return s.getData().SliceDefaultTag
}

func (s *HappyConfig) GetSlices() (map[string]Slice, error) {
	return s.getData().Slices, nil
}
