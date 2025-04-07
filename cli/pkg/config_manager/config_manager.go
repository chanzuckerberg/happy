package config_manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/cli/templates"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Name           string
	ServiceType    string
	Context        string
	DockerfilePath string
	Port           int
	Uri            string
	Priority       int
}

const (
	serviceTypePrivate  = "Service is not exposed to the internet, and can only be consumed by other services in the stack (PRIVATE)"
	serviceTypeExternal = "Service is exposed to the internet (EXTERNAL)"
	serviceTypeInternal = "Service is exposed to the internet, but is protected by OIDC (INTERNAL)"
)

var (
	defaultRegions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}

	serviceTypeMapping = map[string]string{
		serviceTypePrivate:  "PRIVATE",
		serviceTypeExternal: "EXTERNAL",
		serviceTypeInternal: "INTERNAL",
	}

	ErrServiceSkipped = errors.New("service was skipped")
)

type happyConfigDescriptor struct {
	bootstrapConfig *config.Bootstrap
	happyConfig     *config.HappyConfig

	defaultServicePort string
	dockerPaths        []string
	appName            string
	profiles           []string
	environmentNames   []string
	environments       map[string]config.Environment
	services           []Service
}

type configAssembler func(ctx context.Context, descriptor *happyConfigDescriptor) error

func assemble(ctx context.Context, descriptor *happyConfigDescriptor, assemblers ...configAssembler) error {
	for _, assembler := range assemblers {
		err := assembler(ctx, descriptor)
		if err != nil {
			return errors.Wrap(err, "unable to complete a creation step")
		}
	}
	return nil
}

func findServiceCandidates(ctx context.Context, descriptor *happyConfigDescriptor) error {
	dockerPaths, err := findAllDockerfiles(descriptor.bootstrapConfig.HappyProjectRoot)
	if err != nil {
		return errors.Wrap(err, "unable to scan this project for dockerfiles")
	}

	defaultServicePort := "8080"

	logrus.Info("Welcome to happy bootstrap! We'll ask you a few questions to get started.")

	if len(dockerPaths) == 0 {
		logrus.Info("No dockerfiles found in this repo, let us drop one in")
		t, err := templates.StaticAsset("Dockerfile.tmpl")
		if err != nil {
			return errors.Wrap(err, "unable to read a dockerfile template")
		}
		dockerfilePath := filepath.Join(descriptor.bootstrapConfig.HappyProjectRoot, "Dockerfile")
		err = os.WriteFile(dockerfilePath, t, 0644)
		if err != nil {
			return errors.Wrap(err, "unable to create a Dockerfile")
		}
		dockerPaths = append(dockerPaths, dockerfilePath)
		defaultServicePort = "80"
	}
	descriptor.defaultServicePort = defaultServicePort
	descriptor.dockerPaths = dockerPaths

	return nil
}

func appNameExtractor(ctx context.Context, descriptor *happyConfigDescriptor) error {
	appName := ""
	prompt1 := &survey.Input{
		Message: "What would you like to name this application?",
		Help:    "This will be the unique name of the application, lowercased and hyphenated",
		Default: normalizeKey(filepath.Base(descriptor.bootstrapConfig.HappyProjectRoot)),
	}
	err := survey.AskOne(prompt1, &appName,
		survey.WithValidator(survey.Required),
		survey.WithValidator(survey.MinLength(5)),
		survey.WithValidator(survey.MaxLength(15)))
	if err != nil {
		return errors.Wrap(err, "unable to prompt")
	}

	appName = normalizeKey(appName)

	if len(appName) == 0 {
		return errors.New("no application name provided")
	}
	descriptor.appName = appName
	return nil
}

func profileExtractor(ctx context.Context, descriptor *happyConfigDescriptor) error {
	profiles, err := util.GetAWSProfiles()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve aws profiles")
	}
	if len(profiles) == 0 {
		return errors.New("no aws profiles found")
	}
	descriptor.profiles = profiles
	return nil
}

func environmentConfigurator(ctx context.Context, descriptor *happyConfigDescriptor) error {
	environmentNames := []string{}
	prompt := []*survey.Question{
		{
			Name: "environments",
			Prompt: &survey.MultiSelect{
				Message: "Your application will be deployed to multiple environments. Which environments would you like to deploy to?",
				Options: []string{"rdev", "dev", "staging", "prod"},
			},
		},
	}

	err := survey.Ask(prompt, &environmentNames, survey.WithValidator(survey.Required))
	if err != nil {
		return errors.Wrapf(err, "failed to prompt")
	}
	if len(environmentNames) == 0 {
		return errors.New("no environments were selected")
	}

	environments := map[string]config.Environment{}
	for _, env := range environmentNames {
		logrus.Infof("A few questions about environment  '%s':", env)

		environment, err := configureEnvironment(ctx, env, descriptor.profiles)
		if err != nil {
			return errors.Wrapf(err, "failed to configure environment %s", env)
		}

		environments[env] = *environment
	}

	if len(environments) == 0 {
		return errors.New("you have not configured any ehvironments")
	}
	descriptor.environmentNames = environmentNames
	descriptor.environments = environments
	return nil
}

func serviceConfigurator(ctx context.Context, descriptor *happyConfigDescriptor) error {
	services := []Service{}
	logrus.Info("We have found dockerfiles in your project, let's see if you'd like to use them as services in your stack")
	for _, dockerPath := range descriptor.dockerPaths {
		service, err := configureService(descriptor.bootstrapConfig, dockerPath, descriptor.defaultServicePort)

		if errors.Is(err, ErrServiceSkipped) {
			continue
		}

		if err != nil {
			return errors.Wrapf(err, "failed to configure service for dockerfile %s", dockerPath)
		}

		services = append(services, *service)
	}

	if len(services) == 0 {
		return errors.New("you have not configured any services")
	}

	// sort services by length of service.Uri in reverse order
	sort.Slice(services, func(i, j int) bool {
		return len(services[i].Uri) > len(services[j].Uri)
	})
	priority := 0
	for _, service := range services {
		service.Priority = priority
		priority++
	}
	descriptor.services = services
	return nil
}

func CreateHappyConfig(ctx context.Context, bootstrapConfig *config.Bootstrap) (*config.HappyConfig, error) {
	happyConfig, err := config.NewBlankHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get happy config")
	}

	descriptor := &happyConfigDescriptor{
		bootstrapConfig: bootstrapConfig,
		happyConfig:     happyConfig,
	}

	err = assemble(ctx,
		descriptor,
		findServiceCandidates,
		appNameExtractor,
		profileExtractor,
		environmentConfigurator,
		serviceConfigurator)

	if err != nil {
		return nil, errors.Wrap(err, "unable to assemble the app configuration")
	}

	consolidateConfiguration(happyConfig, descriptor)

	err = os.MkdirAll(filepath.Dir(bootstrapConfig.HappyConfigPath), 0777)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create folder %s", filepath.Dir(bootstrapConfig.HappyConfigPath))
	}
	err = happyConfig.Save()

	if err != nil {
		return nil, errors.Wrap(err, "unable to save happy config")
	}

	happyConfig, err = config.NewHappyConfig(bootstrapConfig)

	if err != nil {
		return nil, errors.Wrap(err, "unable to load happy config")
	}

	return happyConfig, nil
}

func consolidateConfiguration(happyConfig *config.HappyConfig, descriptor *happyConfigDescriptor) {
	happyConfig.GetData().Environments = descriptor.environments
	happyConfig.GetData().FeatureFlags = config.Features{
		EnableDynamoLocking:   true,
		EnableHappyApiUsage:   true,
		EnableECRAutoCreation: true,
	}
	happyConfig.GetData().DefaultEnv = descriptor.environmentNames[0]
	happyConfig.GetData().DefaultComposeEnvFile = ".env.ecr"
	happyConfig.GetData().App = descriptor.appName

	serviceDefs := map[string]any{}
	stackDefaults := map[string]any{
		"stack_defaults":   "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main",
		"routing_method":   "CONTEXT",
		"create_dashboard": false,
		"services":         serviceDefs,
	}

	serviceNames := []string{}
	for _, service := range descriptor.services {
		serviceUri := service.Uri
		if !strings.HasSuffix(serviceUri, "/") {
			serviceUri = fmt.Sprintf("%s/", serviceUri)
		}
		serviceDefs[service.Name] =
			map[string]any{
				"name":                  service.Name,
				"port":                  service.Port,
				"health_check_path":     serviceUri,
				"service_type":          service.ServiceType,
				"path":                  fmt.Sprintf("%s*", serviceUri),
				"priority":              service.Priority,
				"success_codes":         "200-499",
				"platform_architecture": "arm64",
				"build": map[string]any{
					"context":    service.Context,
					"dockerfile": service.DockerfilePath,
				},
			}
		serviceNames = append(serviceNames, service.Name)
	}
	happyConfig.GetData().StackDefaults = stackDefaults
	happyConfig.GetData().Services = serviceNames
}
