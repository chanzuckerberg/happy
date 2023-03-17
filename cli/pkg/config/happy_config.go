package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Environment struct {
	AWSProfile         *string         `yaml:"aws_profile"`
	K8S                k8s.K8SConfig   `yaml:"k8s"`
	SecretId           string          `yaml:"secret_arn"`
	TerraformDirectory string          `yaml:"terraform_directory"`
	AutoRunMigrations  bool            `yaml:"auto_run_migrations"`
	TaskLaunchType     util.LaunchType `yaml:"task_launch_type"`
	LogGroupPrefix     string          `yaml:"log_group_prefix"`
}

type Features struct {
	EnableDynamoLocking   bool `yaml:"enable_dynamo_locking"`
	EnableHappyApiUsage   bool `yaml:"enable_happy_api_usage"`
	EnableECRAutoCreation bool `yaml:"enable_ecr_auto_creation"`
}

type ConfigData struct {
	ConfigVersion         string                 `yaml:"config_version"`
	DefaultEnv            string                 `yaml:"default_env"`
	App                   string                 `yaml:"app"`
	DefaultComposeEnvFile string                 `yaml:"default_compose_env_file"`
	Environments          map[string]Environment `yaml:"environments"`
	Tasks                 map[string][]string    `yaml:"tasks"`
	SliceDefaultTag       string                 `yaml:"slice_default_tag"`
	Slices                map[string]Slice       `yaml:"slices"`
	Services              []string               `yaml:"services"`
	FeatureFlags          Features               `yaml:"features"`
	Api                   HappyApiConfig         `yaml:"api"`
}

type Slice struct {
	DeprecatedBuildImages []string `yaml:"build_images"`
	Profile               *Profile `yaml:"profile"`
}

type Profile string

func (p *Profile) Get() string {
	// If no profile specified, we default to "everything"
	// See: https://github.com/docker/compose/issues/8676
	if p == nil {
		return "*"
	}
	return string(*p)
}

type HappyConfig struct {
	env  string
	data *ConfigData

	envConfig *Environment

	projectRoot string
	dockerRepo  string

	composeEnvFile string
}

func NewHappyConfig(bootstrap *Bootstrap) (*HappyConfig, error) {
	configFilePath := bootstrap.GetHappyConfigPath()
	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}

	configData := &ConfigData{}
	err = yaml.Unmarshal(configContent, configData)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing yaml file")
	}

	// validate that DefaultEnv exists in Happy config
	if configData.DefaultEnv == "" {
		return nil, errors.Errorf("Happy config requires a a default environment to be specified under default_env")
	}

	env := bootstrap.GetEnv()
	if len(env) == 0 {
		env = configData.DefaultEnv
	}
	envConfig, ok := configData.Environments[env]
	if !ok {
		return nil, errors.Errorf("environment not found: %s", env)
	}

	// If specified by user, take precedence over config or env
	if bootstrap.GetAWSProfile() != nil {
		envConfig.AWSProfile = bootstrap.GetAWSProfile()
	}

	defaultComposeEnvFile := configData.DefaultComposeEnvFile
	if defaultComposeEnvFile == "" {
		return nil, errors.New("default_compose_env has been superseeded by default_compose_env_file")
	}

	composeEnvFile, err := findDockerComposeFile(bootstrap)
	if err != nil {
		return nil, err
	}

	happyRootPath := bootstrap.GetHappyProjectRootPath()

	config := &HappyConfig{
		env:            env,
		data:           configData,
		envConfig:      &envConfig,
		composeEnvFile: composeEnvFile,

		projectRoot: happyRootPath,
	}

	return config, config.validate()
}

// validate validates the config
func (s *HappyConfig) validate() error {
	// TODO: there is probably a bunch of other validation we need
	var deprecated error
	for name, slice := range s.data.Slices {
		if len(slice.DeprecatedBuildImages) > 0 {
			deprecated = multierror.Append(
				deprecated,
				errors.Errorf(
					"slice(%s).build_images is deprecated and will be ignored. please use profiles instead.",
					name,
				),
			)
		}
	}
	if deprecated != nil {
		logrus.Debug(deprecated)
	}

	return nil
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

func (s *HappyConfig) AwsProfile() *string {
	envConfig := s.getEnvConfig()

	return envConfig.AWSProfile
}

func (s *HappyConfig) GetSecretId() string {
	envConfig := s.getEnvConfig()

	return envConfig.SecretId
}

func (s *HappyConfig) GetLogGroupPrefix() string {
	envConfig := s.getEnvConfig()

	return envConfig.LogGroupPrefix
}

func (s *HappyConfig) AutoRunMigrations() bool {
	envConfig := s.getEnvConfig()

	return envConfig.AutoRunMigrations
}

func (s *HappyConfig) TerraformDirectory() string {
	envConfig := s.getEnvConfig()

	return envConfig.TerraformDirectory
}

func (s *HappyConfig) TaskLaunchType() util.LaunchType {
	envConfig := s.getEnvConfig()

	taskLaunchType := util.LaunchType(strings.ToUpper(envConfig.TaskLaunchType.String()))
	if taskLaunchType != util.LaunchTypeFargate && taskLaunchType != util.LaunchTypeK8S {
		taskLaunchType = util.LaunchTypeEC2
	}
	return taskLaunchType
}

func (s *HappyConfig) K8SConfig() *k8s.K8SConfig {
	envConfig := s.getEnvConfig()
	return &envConfig.K8S
}

func (s *HappyConfig) DefaultEnv() string {
	return s.getData().DefaultEnv
}

func (s *HappyConfig) DefaultComposeEnvFile() string {
	return s.getData().DefaultComposeEnvFile
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

func (s *HappyConfig) TaskExists(taskType string) bool {
	_, ok := s.getData().Tasks[taskType]
	return ok
}

func (s *HappyConfig) GetServices() []string {
	return s.getData().Services
}

func (s *HappyConfig) SliceDefaultTag() string {
	return s.getData().SliceDefaultTag
}

func (s *HappyConfig) GetSlice(name string) (*Slice, error) {
	slices := s.getData().Slices
	slice, found := slices[name]
	if !found {
		return nil, errors.Errorf("slice(%s) is not a valid slice.", name)
	}
	return &slice, nil
}

func (s *HappyConfig) GetDockerRepo() string {
	return s.dockerRepo
}

func (s *HappyConfig) GetDockerComposeEnvFile() string {
	return s.composeEnvFile
}

func (s *HappyConfig) GetFeatures() *Features {
	return &s.getData().FeatureFlags
}

func (s *HappyConfig) GetHappyApiConfig() HappyApiConfig {
	defaultApiConfig := DefaultHappyApiConfig()

	apiConfig := s.getData().Api
	if apiConfig.BaseUrl == "" {
		apiConfig.BaseUrl = defaultApiConfig.BaseUrl
	}
	if apiConfig.OidcClientID == "" {
		apiConfig.OidcClientID = defaultApiConfig.OidcClientID
	}
	if apiConfig.OidcIssuerUrl == "" {
		apiConfig.OidcIssuerUrl = defaultApiConfig.OidcIssuerUrl
	}
	return apiConfig
}

func findDockerComposeFile(bootstrap *Bootstrap) (string, error) {
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
	absComposeEnvFile, err := findFile(composeEnvFile, pathsToLook)
	if err != nil {
		return "", errors.Wrapf(err, "cannot locate docker-compose env file %s", composeEnvFile)
	}
	return absComposeEnvFile, nil
}

func findFile(fileName string, paths []string) (string, error) {
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
			logrus.Infof("Looking for %s\n", filePath)
			file, err := os.Stat(filePath)
			if err == nil && !file.IsDir() {
				return filePath, nil
			}
		}

		return "", errors.Errorf("cannot locate file %s anywhere", fileName)
	}
}
