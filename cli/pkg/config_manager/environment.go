package config_manager

import (
	"context"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func configureEnvironment(ctx context.Context, env string, profiles []string) (*config.Environment, error) {
	profile := ""
	region := ""
	clusterId := ""
	happyNamespace := ""

	for {
		prompt := &survey.Select{
			Message: fmt.Sprintf("Which aws profile do you want to use in %s?", env),
			Options: profiles,
			Default: profiles[0],
		}

		err := survey.AskOne(prompt, &profile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prompt")
		}
		if len(profile) == 0 {
			continue
		}

		var clusterIds []string
		prompt1 := &survey.Select{
			Message: fmt.Sprintf("Which aws region should we use in %s? (us-west-1 should be avoided, if possible)", env),
			Options: defaultRegions,
			Default: "us-west-2",
		}

		err = survey.AskOne(prompt1, &region)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prompt")
		}
		logrus.Info("Checking for eks clusters in this region...")

		clusterIds, err = listEKSClusterIDs(ctx, profile, region)
		if err != nil {
			return nil, errors.Wrap(err, "unable to list eks clusters")
		}
		if len(clusterIds) == 0 {
			logrus.Error("No eks clusters found in this region. Please select a different region or profile.")
			continue
		}

		prompt = &survey.Select{
			Message: fmt.Sprintf("Which EKS cluster should we use in %s?", env),
			Options: clusterIds,
			Default: clusterIds[0],
		}

		err = survey.AskOne(prompt, &clusterId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prompt")
		}

		logrus.Info("Checking for happy environments in this cluster...")

		var happyNamespaces []string
		happyNamespaces, err = listHappyNamespaces(ctx, profile, region, clusterId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to obtain a list of happy namespaces")
		}
		if len(happyNamespaces) == 0 {
			logrus.Error("No happy namespaces were found in the selected cluster, please select a different region, profile or eks cluster.")
			continue
		}

		defaultNamespace := getDefaultHappyNamespace(happyNamespaces, env)

		prompt2 := &survey.Select{
			Message: fmt.Sprintf("Which happy namespace should we use in %s?", env),
			Options: happyNamespaces,
			Default: defaultNamespace,
		}

		err = survey.AskOne(prompt2, &happyNamespace)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to obtain an eks cluster id")
		}
		break
	}

	return &config.Environment{
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
	}, nil
}
