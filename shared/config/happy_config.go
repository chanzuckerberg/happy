package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_HAPPY_API_BASE_URL        = "https://hapi.hapi.prod.si.czi.technology"
	DEFAULT_HAPPY_API_OIDC_CLIENT_ID  = "0oa8anwuhpAX1rfvb5d7"
	DEFAULT_HAPPY_API_OIDC_ISSUER_URL = "https://czi-prod.okta.com"
)

type Environment struct {
	AWSProfile         *string         `yaml:"aws_profile"`
	AWSRegion          *string         `yaml:"aws_region" default:"us-west-2"`
	K8S                k8s.K8SConfig   `yaml:"k8s"`
	SecretId           string          `yaml:"secret_arn"`
	TerraformDirectory string          `yaml:"terraform_directory"`
	AutoRunMigrations  bool            `yaml:"auto_run_migrations"`
	TaskLaunchType     util.LaunchType `yaml:"task_launch_type"`
	LogGroupPrefix     string          `yaml:"log_group_prefix"`
	StackOverrides     StackConfig     `yaml:"stack_overrides"`
}

type EnvironmentContext struct {
	EnvironmentName string
	AWSProfile      *string
	AWSRegion       *string
	K8S             k8s.K8SConfig
	SecretId        string
	TaskLaunchType  util.LaunchType
}

type Features struct {
	EnableDynamoLocking   bool `yaml:"enable_dynamo_locking"`
	EnableHappyApiUsage   bool `yaml:"enable_happy_api_usage"`
	EnableECRAutoCreation bool `yaml:"enable_ecr_auto_creation"`
}

type HappyApiConfig struct {
	BaseUrl       string `yaml:"base_url"`
	OidcClientID  string `yaml:"oidc_client_id"`
	OidcIssuerUrl string `yaml:"oidc_issuer_url"`
}

type ServiceConfig struct {
	Name                          *string `yaml:"name" json:"name"`
	DesiredCount                  *int    `yaml:"desired_count" json:"desired_count"`
	MaxCount                      *int    `yaml:"max_count" json:"max_count"`
	ScalingCPUThresholdPercentage *int    `yaml:"scaling_cpu_threshold_percentage" json:"scaling_cpu_threshold_percentage"`
	Port                          *int    `yaml:"port" json:"port"`
	Memory                        *string `yaml:"memory" json:"memory"`
	CPU                           *string `yaml:"cpu" json:"cpu"`
	HealthCheckPath               *string `yaml:"health_check_path" json:"health_check_path"`
	ServiceType                   *string `yaml:"service_type" json:"service_type"`
	Path                          *string `yaml:"path" json:"path"`
	Priority                      *int    `yaml:"priority" json:"priority"`
	SuccessCodes                  *string `yaml:"success_codes" json:"success_codes"`
	InitialDelaySeconds           *int    `yaml:"initial_delay_seconds" json:"initial_delay_seconds"`
	PeriodSeconds                 *int    `yaml:"period_seconds" json:"period_seconds"`
	PlatformArchitecture          *string `yaml:"platform_architecture" json:"platform_architecture"`
	AwsIamPolicyJSON              *string `yaml:"aws_iam_policy_json" json:"aws_iam_policy_json"`
	Synthetics                    *bool   `yaml:"synthetics" json:"synthetics"`
}

type StackConfig struct {
	Source          *string                  `yaml:"source" json:"source"`
	RoutingMethod   *string                  `yaml:"routing_method" json:"routing_method"`
	CreateDashboard *bool                    `yaml:"create_dashboard" json:"create_dashboard"`
	Services        map[string]ServiceConfig `yaml:"services" json:"services"`
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
	StackDefaults         StackConfig            `yaml:"stack_defaults"`
}

type Slice struct {
	DeprecatedBuildImages []string `yaml:"build_images"`
	Profile               *Profile `yaml:"profile"`
}

func (ec *EnvironmentContext) GetEnv() string {
	return ec.EnvironmentName
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
	err = defaults.Set(configData)
	if err != nil {
		return nil, errors.Wrap(err, "error setting config defaults")
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

func (s *HappyConfig) GetEnvConfig() *Environment {
	return s.envConfig
}

func (s *HappyConfig) GetEnv() string {
	return s.env
}

func (s *HappyConfig) GetProjectRoot() string {
	return s.projectRoot
}

func (s *HappyConfig) AwsProfile() *string {
	envConfig := s.GetEnvConfig()

	return envConfig.AWSProfile
}

func (s *HappyConfig) AwsRegion() *string {
	envConfig := s.GetEnvConfig()

	return envConfig.AWSRegion
}

func (s *HappyConfig) GetSecretId() string {
	envConfig := s.GetEnvConfig()

	return envConfig.SecretId
}

func (s *HappyConfig) GetLogGroupPrefix() string {
	envConfig := s.GetEnvConfig()

	return envConfig.LogGroupPrefix
}

func (s *HappyConfig) AutoRunMigrations() bool {
	envConfig := s.GetEnvConfig()

	return envConfig.AutoRunMigrations
}

func (s *HappyConfig) TerraformDirectory() string {
	envConfig := s.GetEnvConfig()

	return envConfig.TerraformDirectory
}

func (s *HappyConfig) TaskLaunchType() util.LaunchType {
	envConfig := s.GetEnvConfig()

	taskLaunchType := util.LaunchType(strings.ToUpper(envConfig.TaskLaunchType.String()))
	if taskLaunchType != util.LaunchTypeFargate && taskLaunchType != util.LaunchTypeK8S {
		taskLaunchType = util.LaunchTypeEC2
	}
	return taskLaunchType
}

// Recursively combines stack_defaults and stack_overrides if they are set; returns structured config, and unstructured config
func (s *HappyConfig) GetStackConfig() (*StackConfig, map[string]interface{}, error) {
	src := map[string]interface{}{}
	dst := map[string]interface{}{}
	stackConfig := &StackConfig{}

	// Convert stack configuration to a nested map
	err := util.DeepClone(&dst, s.getData().StackDefaults)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot serialize stack defaults")
	}
	err = util.DeepClone(&src, s.GetEnvConfig().StackOverrides)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot serialize stack overrides")
	}

	// merge two maps together, recursively (ignoring null values)
	util.DeepMerge(dst, src)

	// Convert nested map back to StackConfig struct
	err = util.DeepClone(stackConfig, dst)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot deserialize merged stack config")
	}

	return stackConfig, dst, nil
}

func (s *HappyConfig) K8SConfig() *k8s.K8SConfig {
	envConfig := s.GetEnvConfig()
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

func (s *HappyConfig) GetStackDefaults() *StackConfig {
	return &s.getData().StackDefaults
}

func (s *HappyConfig) GetHappyApiConfig() HappyApiConfig {
	apiConfig := s.getData().Api
	if apiConfig.BaseUrl == "" {
		apiConfig.BaseUrl = DEFAULT_HAPPY_API_BASE_URL
	}
	if apiConfig.OidcClientID == "" {
		apiConfig.OidcClientID = DEFAULT_HAPPY_API_OIDC_CLIENT_ID
	}
	if apiConfig.OidcIssuerUrl == "" {
		apiConfig.OidcIssuerUrl = DEFAULT_HAPPY_API_OIDC_ISSUER_URL
	}
	return apiConfig
}

func (s *HappyConfig) GetEnvironmentContext() EnvironmentContext {
	return EnvironmentContext{
		EnvironmentName: s.GetEnv(),
		AWSProfile:      s.AwsProfile(),
		AWSRegion:       s.AwsRegion(),
		K8S:             *s.K8SConfig(),
		SecretId:        s.GetSecretId(),
		TaskLaunchType:  s.TaskLaunchType(),
	}
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
