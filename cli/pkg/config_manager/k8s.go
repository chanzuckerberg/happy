package config_manager

import (
	"context"
	"fmt"
	"strings"

	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
)

func listEKSClusterIDs(ctx context.Context, profile, region string) ([]string, error) {
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

func listHappyNamespaces(ctx context.Context, profile, region, clusterId string) ([]string, error) {
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

func getDefaultHappyNamespace(happyNamespaces []string, env string) string {
	defaultNamespace := happyNamespaces[0]
	for _, namespace := range happyNamespaces {
		if strings.Contains(namespace, fmt.Sprintf("-%s-happy-env", env)) {
			defaultNamespace = namespace
			break
		}
	}
	return defaultNamespace
}
