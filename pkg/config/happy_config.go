package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	ConfigVersion         string                 `yaml:"config_version"`
	TerraformVersion      string                 `yaml:"terraform_version"`
	DefaultEnv            string                 `yaml:"default_env"`
	App                   string                 `yaml:"app"`
	DefaultComposeEnvFile string                 `yaml:"default_compose_env_file"`
	Environments          map[string]Environment `yaml:"environments"`
	Tasks                 map[string][]string    `yaml:"tasks"`
	SliceDefaultTag       string                 `yaml:"slice_default_tag"`
	Slices                map[string]Slice       `yaml:"slices"`
	Services              []string               `yaml:"services"`
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
	GetEnv() string
	DefaultComposeEnvFile() string
	App() string
	GetServiceRegistries() map[string]*RegistryConfig
	ClusterArn() string
	PrivateSubnets() []string
	SecurityGroups() []string
	TfeUrl() string
	TfeOrg() string
	SliceDefaultTag() string
	GetSlices() (map[string]Slice, error)
	TaskLaunchType() string
	GetServices() []string
	GetDockerRepo() string
	GetComposeEnvFile() string
}

type happyConfig struct {
	env  string
	data *ConfigData

	envConfig *Environment

	projectRoot string
	dockerRepo  string

	serviceRegistries map[string]*RegistryConfig
	clusterArn        string
	privateSubnets    []string
	securityGroups    []string
	tfeUrl            string
	tfeOrg            string

	composeEnvFile string
}

func NewHappyConfig(bootstrap *Bootstrap) (HappyConfig, error) {
	return NewHappyConfigWithSecretsBackend(bootstrap, nil)
}

func NewHappyConfigWithSecretsBackend(bootstrap *Bootstrap, secretMgr SecretsBackend) (HappyConfig, error) {
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
	if len(env) == 0 {
		env = configData.DefaultEnv
	}
	envConfig, ok := configData.Environments[env]
	if !ok {
		return nil, errors.Errorf("environment not found: %s", env)
	}

	defaultComposeEnvFile := configData.DefaultComposeEnvFile
	if len(defaultComposeEnvFile) == 0 {
		return nil, errors.New("default_compose_env has been superseeded by default_compose_env_file")
	}

	composeEnvFile := bootstrap.GetComposeEnvFile()
	if len(composeEnvFile) == 0 {
		composeEnvFile = defaultComposeEnvFile
	}

	if len(composeEnvFile) > 0 && len(bootstrap.GetHappyProjectRootPath()) > 0 {
		absComposeEnvFile, err := FindDockerComposeFile(bootstrap, composeEnvFile)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot locate docker-compose env file %s", composeEnvFile)
		}
		composeEnvFile = absComposeEnvFile
	}

	awsProfile := envConfig.AWSProfile
	if secretMgr == nil {
		secretMgr = GetAwsSecretMgr(awsProfile)
	}
	secretArn := envConfig.SecretARN

	secrets, err := secretMgr.GetSecrets(secretArn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve secrets")
	}

	dockerRepo := os.Getenv("DOCKER_REPO")
	if len(dockerRepo) == 0 {
		serviceRegistries := secrets.GetAllServicesUrl()
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

	happyRootPath := bootstrap.GetHappyProjectRootPath()

	config := &happyConfig{
		env:       env,
		data:      &configData,
		envConfig: &envConfig,

		projectRoot: happyRootPath,

		serviceRegistries: secrets.GetAllServicesUrl(),
		dockerRepo:        dockerRepo,
		clusterArn:        secrets.GetClusterArn(),
		privateSubnets:    secrets.GetPrivateSubnets(),
		securityGroups:    secrets.GetSecurityGroups(),
		tfeUrl:            secrets.GetTfeUrl(),
		tfeOrg:            secrets.GetTfeOrg(),
		composeEnvFile:    composeEnvFile,
	}

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

func (s *happyConfig) DefaultComposeEnvFile() string {
	return s.getData().DefaultComposeEnvFile
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

func (s *happyConfig) GetServiceRegistries() map[string]*RegistryConfig {
	return s.serviceRegistries
}

func (s *happyConfig) ClusterArn() string {
	return s.clusterArn
}

func (s *happyConfig) PrivateSubnets() []string {
	return s.privateSubnets
}

func (s *happyConfig) SecurityGroups() []string {
	return s.securityGroups
}

func (s *happyConfig) TfeUrl() string {
	return s.tfeUrl
}

func (s *happyConfig) TfeOrg() string {
	return s.tfeOrg
}

func (s *happyConfig) SliceDefaultTag() string {
	return s.getData().SliceDefaultTag
}

func (s *happyConfig) GetSlices() (map[string]Slice, error) {
	return s.getData().Slices, nil
}

func (s *happyConfig) GetDockerRepo() string {
	return s.dockerRepo
}

func (s *happyConfig) GetComposeEnvFile() string {
	return s.composeEnvFile
}

func FindDockerComposeFile(bootstrap *Bootstrap, fileName string) (string, error) {
	if bootstrap == nil || len(fileName) == 0 {
		return fileName, nil
	}

	// Look in the project root first, then current directory, then home directory, then parent directory, then parent of a parent directory
	pathsToLook := []string{bootstrap.GetHappyProjectRootPath()}
	currentDir, err := os.Getwd()
	if err == nil {
		pathsToLook = append(pathsToLook, currentDir)
	}
	homeDir, err := os.UserHomeDir()
	if err == nil {
		pathsToLook = append(pathsToLook, homeDir)
	}
	parentDir, err := filepath.Abs("..")
	if err == nil {
		pathsToLook = append(pathsToLook, parentDir)
	}
	grandParentDir, err := filepath.Abs("../..")
	if err == nil {
		pathsToLook = append(pathsToLook, grandParentDir)
	}
	absComposeEnvFile, err := FindFile(composeEnvFile, pathsToLook)
	if err != nil {
		return "", errors.Wrapf(err, "cannot locate docker-compose env file %s", composeEnvFile)
	}
	return absComposeEnvFile, nil
}

func FindFile(fileName string, paths []string) (string, error) {
	if len(fileName) == 0 {
		return fileName, nil
	}
	if filepath.IsAbs(fileName) {
		file, err := os.Stat(fileName)
		if err != nil {
			return "", errors.Wrap(err, "cannot find file")
		}
		if file.IsDir() {
			return "", errors.Errorf("provided path %s is a directory", fileName)
		}
		return fileName, nil
	} else {
		for _, path := range paths {
			filePath := filepath.Join(path, fileName)
			log.Printf("Looking for %s\n", filePath)
			file, err := os.Stat(filePath)
			if err == nil && !file.IsDir() {
				return filePath, nil
			}
		}

		return "", errors.Errorf("cannot locate file %s anywhere", fileName)
	}
}
