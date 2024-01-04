package config

import (
	"encoding/json"
	"fmt"
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
	AWSProfile         *string         `yaml:"aws_profile" json:"aws_profile,omitempty"`
	AWSRegion          *string         `yaml:"aws_region" json:"aws_region,omitempty" default:"us-west-2"`
	K8S                k8s.K8SConfig   `yaml:"k8s" json:"k8s,omitempty"`
	SecretId           string          `yaml:"secret_arn" json:"secret_arn,omitempty"`
	TerraformDirectory string          `yaml:"terraform_directory" json:"terraform_directory,omitempty"`
	AutoRunMigrations  bool            `yaml:"auto_run_migrations" json:"auto_run_migrations,omitempty"`
	TaskLaunchType     util.LaunchType `yaml:"task_launch_type" json:"task_launch_type,omitempty"`
	LogGroupPrefix     string          `yaml:"log_group_prefix" json:"log_group_prefix,omitempty"`
	StackOverrides     map[string]any  `yaml:"stack_overrides" json:"stack_overrides,omitempty"`
}

type EnvironmentContext struct {
	EnvironmentName string
	AWSProfile      *string
	AWSRegion       *string
	K8S             k8s.K8SConfig
	SecretID        string
	TaskLaunchType  util.LaunchType
}

type Features struct {
	EnableDynamoLocking       bool `yaml:"enable_dynamo_locking" json:"enable_dynamo_locking,omitempty"`
	EnableHappyApiUsage       bool `yaml:"enable_happy_api_usage" json:"enable_happy_api_usage,omitempty"`
	EnableECRAutoCreation     bool `yaml:"enable_ecr_auto_creation" json:"enable_ecr_auto_creation,omitempty"`
	EnableUnifiedConfig       bool `yaml:"enable_unified_config" json:"enable_unified_config,omitempty"`
	EnableUnusedImageDeletion bool `yaml:"enable_unused_image_deletion" json:"enable_unused_image_deletion,omitempty"`
	EnableHappyConfigV2       bool `yaml:"enable_happy_config_v2" json:"enable_happy_config_v2,omitempty"`
}

type HappyApiConfig struct {
	BaseUrl       string `yaml:"base_url" json:"base_url,omitempty"`
	OidcClientID  string `yaml:"oidc_client_id" json:"oidc_client_id,omitempty"`
	OidcIssuerUrl string `yaml:"oidc_issuer_url" json:"oidc_issuer_url,omitempty"`
}

type ConfigData struct {
	ConfigVersion         string                 `yaml:"config_version" json:"config_version,omitempty"`
	DefaultEnv            string                 `yaml:"default_env" json:"default_env,omitempty"`
	App                   string                 `yaml:"app" json:"app,omitempty"`
	DefaultComposeEnvFile string                 `yaml:"default_compose_env_file" json:"default_compose_env_file,omitempty"`
	Environments          map[string]Environment `yaml:"environments" json:"environments,omitempty"`
	Tasks                 map[string][]string    `yaml:"tasks" json:"tasks,omitempty"`
	SliceDefaultTag       string                 `yaml:"slice_default_tag" json:"slice_default_tag,omitempty"`
	Slices                map[string]Slice       `yaml:"slices" json:"slices,omitempty"`
	Services              []string               `yaml:"services" json:"services,omitempty"`
	FeatureFlags          Features               `yaml:"features" json:"features,omitempty"`
	Api                   HappyApiConfig         `yaml:"api" json:"api,omitempty"`
	StackDefaults         map[string]any         `yaml:"stack_defaults" json:"stack_defaults,omitempty"`
}

type Slice struct {
	DeprecatedBuildImages []string `yaml:"build_images" json:"build_images,omitempty"`
	Profile               *Profile `yaml:"profile" json:"profile,omitempty"`
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
	bootstrap *Bootstrap

	projectRoot string
	dockerRepo  string

	composeEnvFile string
	configFilePath string
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

	if composeEnvFile == "" {
		composeEnvFile = defaultComposeEnvFile
	}

	absComposeEnvFile, err := findDockerComposeEnvFile(bootstrap)
	if err != nil {
		logrus.Debugf("Unable to find docker-compose env file %s: %s", composeEnvFile, err.Error())
	}

	happyRootPath := bootstrap.GetHappyProjectRootPath()

	config := &HappyConfig{
		env:            env,
		data:           configData,
		bootstrap:      bootstrap,
		envConfig:      &envConfig,
		composeEnvFile: absComposeEnvFile,

		projectRoot:    happyRootPath,
		configFilePath: configFilePath,
	}

	return config, config.validate()
}

func NewBlankHappyConfig(bootstrap *Bootstrap) (*HappyConfig, error) {
	configFilePath := bootstrap.GetHappyConfigPath()
	configData := &ConfigData{
		Environments: map[string]Environment{},
	}

	env := bootstrap.GetEnv()
	if len(env) == 0 {
		env = configData.DefaultEnv
	}

	envConfig := Environment{}
	configData.Environments[env] = envConfig

	happyRootPath := bootstrap.GetHappyProjectRootPath()

	config := &HappyConfig{
		env:            env,
		data:           configData,
		bootstrap:      bootstrap,
		envConfig:      &envConfig,
		composeEnvFile: composeEnvFile,

		projectRoot:    happyRootPath,
		configFilePath: configFilePath,
	}

	return config, nil
}

func (s *HappyConfig) Save() error {
	d, err := json.MarshalIndent(s.GetData(), "", "    ")
	if err != nil {
		return errors.Wrapf(err, "Unable to convert config struct to json")
	}
	err = os.WriteFile(s.configFilePath, d, 0644)
	return errors.Wrap(err, "Unable to write config file")
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

func (s *HappyConfig) GetData() *ConfigData {
	return s.data
}

func (s *HappyConfig) SetStackDefaults(stackDefaults map[string]any) {
	s.GetData().StackDefaults = stackDefaults
}

func (s *HappyConfig) GetBootstrap() *Bootstrap {
	return s.bootstrap
}

func (s *HappyConfig) GetEnvConfig() *Environment {
	return s.envConfig
}

func (s *HappyConfig) GetEnv() string {
	return s.env
}

// Only to be used for bootstrapping in happy bootstrap.
func (s *HappyConfig) SetEnv(env string) {
	s.env = env
	envConfig := s.data.Environments[env]
	s.envConfig = &envConfig
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
	if taskLaunchType != util.LaunchTypeFargate && taskLaunchType != util.LaunchTypeEC2 && taskLaunchType != util.LaunchTypeK8S {
		taskLaunchType = util.LaunchTypeK8S
	}
	return taskLaunchType
}

// Recursively combines stack_defaults and stack_overrides if they are set; returns structured config, and unstructured config
func (s *HappyConfig) GetStackConfig() (map[string]interface{}, error) {
	src := map[string]interface{}{}
	dst := map[string]interface{}{}

	// Convert stack configuration to a nested map
	err := util.DeepClone(&dst, s.GetData().StackDefaults)
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize stack defaults")
	}
	err = util.DeepClone(&src, s.GetEnvConfig().StackOverrides)
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize stack overrides")
	}

	// merge two maps together, recursively (ignoring null values)
	err = util.DeepMerge(dst, src)
	if err != nil {
		return nil, errors.Wrap(err, "cannot merge stack defaults and overrides")
	}

	return dst, nil
}

func (s *HappyConfig) K8SConfig() *k8s.K8SConfig {
	envConfig := s.GetEnvConfig()
	return &envConfig.K8S
}

func (s *HappyConfig) DefaultEnv() string {
	return s.GetData().DefaultEnv
}

func (s *HappyConfig) DefaultComposeEnvFile() string {
	return s.GetData().DefaultComposeEnvFile
}

func (s *HappyConfig) App() string {
	return s.GetData().App
}

func (s *HappyConfig) GetTasks(taskType string) ([]string, error) {
	tasks, ok := s.GetData().Tasks[taskType]
	if !ok {
		return nil, errors.Errorf("failed to get tasks: task type not found: %s", taskType)
	}
	return tasks, nil
}

func (s *HappyConfig) TaskExists(taskType string) bool {
	_, ok := s.GetData().Tasks[taskType]
	return ok
}

func (s *HappyConfig) GetServices() []string {
	return s.GetData().Services
}

func (s *HappyConfig) SliceDefaultTag() string {
	return s.GetData().SliceDefaultTag
}

func (s *HappyConfig) GetSlice(name string) (*Slice, error) {
	slices := s.GetData().Slices
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
	return &s.GetData().FeatureFlags
}

func (s *HappyConfig) GetHappyAPIConfig() HappyApiConfig {
	apiConfig := s.GetData().Api
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
		SecretID:        s.GetSecretId(),
		TaskLaunchType:  s.TaskLaunchType(),
	}
}

func (s *HappyConfig) GetModuleSource() string {
	moduleSource := ""
	if overrideSource, ok := s.GetEnvConfig().StackOverrides["source"]; ok {
		moduleSource = overrideSource.(string)
	}

	if len(moduleSource) == 0 {
		if defaultSource, ok := s.GetData().StackDefaults["source"]; ok {
			moduleSource = defaultSource.(string)
		}
	}

	if len(moduleSource) == 0 {
		computeFlavor := "eks"
		moduleSource = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-%s?ref=main"
		if s.TaskLaunchType() == util.LaunchTypeFargate || s.TaskLaunchType() == util.LaunchTypeEC2 {
			computeFlavor = "ecs"
		}
		moduleSource = fmt.Sprintf(moduleSource, computeFlavor)
	}

	return moduleSource
}

func (s *HappyConfig) GetModuleNames() map[string]bool {
	if s.TaskLaunchType() == util.LaunchTypeK8S {
		return map[string]bool{"happy-stack-eks": true, "happy-stack-helm-eks": true}
	} else {
		return map[string]bool{"happy-stack-ecs": true}
	}
}

func findDockerComposeEnvFile(bootstrap *Bootstrap) (string, error) {
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
			logrus.Debugf("Looking for %s\n", filePath)
			file, err := os.Stat(filePath)
			if err == nil && !file.IsDir() {
				return filePath, nil
			}
		}

		return "", errors.Errorf("cannot locate file %s anywhere", fileName)
	}
}
