package config_manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func CreeateHappyConfig(ctx context.Context, bootstrapConfig *config.Bootstrap) (*config.HappyConfig, error) {
	logrus.Infof("Bootstrap config: %+v", bootstrapConfig)

	happyConfig, err := config.NewBlankHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get happy config")
	}

	profiles, err := util.GetAwsProfiles()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve aws profiles")
	}
	if len(profiles) == 0 {
		return nil, errors.New("no aws profiles found")
	}

	environmentNames := []string{}
	prompt := []*survey.Question{
		{
			Name: "environments",
			Prompt: &survey.MultiSelect{
				Message: "Which environments should we create?",
				Options: []string{"rdev", "staging", "prod"},
			},
		},
	}

	err = survey.Ask(prompt, &environmentNames)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain an eks cluster id")
	}
	if len(environmentNames) == 0 {
		return nil, errors.New("no environments were selected")
	}

	environments := map[string]config.Environment{}
	for _, env := range environmentNames {
		logrus.Infof("A few questions about %s environment", env)
		profile := ""
		prompt := &survey.Select{
			Message: fmt.Sprintf("Which aws profile do you want to use in %s?", env),
			Options: profiles,
			Default: profiles[0],
		}

		err = survey.AskOne(prompt, &profile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to obtain an aws profile")
		}

		region := ""
		prompt = &survey.Select{
			Message: fmt.Sprintf("Which aws region should we use in %s?", env),
			Options: []string{"us-west-2", "us-east-1", "us-east-2", "us-west-1"},
			Default: "us-west-2",
		}

		err = survey.AskOne(prompt, &region)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to obtain an aws profile")
		}

		clusterIds, err := ListClusterIds(ctx, profile, region)
		if err != nil {
			return nil, errors.Wrap(err, "unable to list eks clusters")
		}
		if len(clusterIds) == 0 {
			return nil, errors.New("no eks clusters found")
		}

		clusterId := ""
		prompt = &survey.Select{
			Message: fmt.Sprintf("Which EKS cluster should we use in %s?", env),
			Options: clusterIds,
			Default: clusterIds[0],
		}

		err = survey.AskOne(prompt, &clusterId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to obtain an eks cluster id")
		}

		happyNamespaces, err := ListHappyNamespaces(ctx, profile, region, clusterId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to obtain a list of happy namespaces")
		}
		if len(happyNamespaces) == 0 {
			return nil, errors.New("no happy namespaces found in the cluster")
		}

		defaultNamespace := happyNamespaces[0]
		for _, namespace := range happyNamespaces {
			if strings.Contains(namespace, fmt.Sprintf("-%s-happy-env", env)) {
				defaultNamespace = namespace
				break
			}
		}

		happyNamespace := ""
		prompt = &survey.Select{
			Message: fmt.Sprintf("Which happy namespace should we use in %s?", env),
			Options: happyNamespaces,
			Default: defaultNamespace,
		}

		err = survey.AskOne(prompt, &happyNamespace)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to obtain an eks cluster id")
		}
		environments[env] = config.Environment{
			AWSProfile: util.String(profile),
			AWSRegion:  util.String(region),
			K8S: k8s.K8SConfig{
				Namespace:  happyNamespace,
				ClusterID:  clusterId,
				AuthMethod: "eks",
			},
			TerraformDirectory: fmt.Sprintf(".happy/terraform/envs/%s", env),
			AutoRunMigrations:  false,
			TaskLaunchType:     "k8s",
		}
	}
	happyConfig.GetData().Environments = environments
	happyConfig.GetData().FeatureFlags = config.Features{
		EnableDynamoLocking:   true,
		EnableHappyApiUsage:   true,
		EnableECRAutoCreation: true,
	}
	happyConfig.GetData().DefaultEnv = environmentNames[0]
	happyConfig.GetData().DefaultComposeEnvFile = ".env.ecr"
	happyConfig.GetData().StackDefaults = map[string]any{
		"stack_defaults": "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main",
		"routing_method": "CONTEXT",

		"create_dashboard": false,
		"services": map[string]any{
			"frontend": map[string]any{
				"name":                             "frontend",
				"desired_count":                    1,
				"max_count":                        1,
				"scaling_cpu_threshold_percentage": 80,
				"port":                             3000,
				"memory":                           "128Mi",
				"cpu":                              "100m",
				"health_check_path":                "/",
				"service_type":                     "INTERNAL",
				"path":                             "/*",
				"priority":                         0,
				"success_codes":                    "200-499",
				"initial_delay_seconds":            30,
				"period_seconds":                   3,
				"platform_architecture":            "arm64",
				"synthetics":                       false,
			},
		},
	}
	happyConfig.GetData().Services = []string{"frontend"}

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

func ListClusterIds(ctx context.Context, profile, region string) ([]string, error) {
	b, err := backend.NewAWSBackend(ctx, config.EnvironmentContext{
		AWSProfile:     util.String(profile),
		AWSRegion:      util.String(region),
		TaskLaunchType: util.LaunchTypeNull,
	})
	if err != nil {
		return []string{}, errors.Wrap(err, "unable to create an aws backend")
	}
	return b.ListEKSClusterIds(ctx)
}

func ListHappyNamespaces(ctx context.Context, profile, region, clusterId string) ([]string, error) {
	b, err := backend.NewAWSBackend(ctx, config.EnvironmentContext{
		AWSProfile:     util.String(profile),
		AWSRegion:      util.String(region),
		TaskLaunchType: util.LaunchTypeK8S,
		K8S: k8s.K8SConfig{
			AuthMethod: "eks",
			ClusterID:  clusterId,
		},
	}, backend.WithIntegrationSecret(&config.IntegrationSecret{})) // This will prevent the backend from trying to load the integration secret, as we have not selected the namespace yet
	if err != nil {
		return []string{}, errors.Wrap(err, "unable to create an aws backend")
	}
	return b.ListHappyNamespaces(ctx)
}
