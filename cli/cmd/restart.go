package cmd

import (
	"context"
	"sync"
	"time"

	happyCmd "github.com/chanzuckerberg/happy/cli/pkg/cmd"
	"github.com/chanzuckerberg/happy/cli/pkg/hapi"
	"github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	rootCmd.AddCommand(restartCmd)
	config.ConfigureCmdWithBootstrapConfig(restartCmd)
}

var restartCmd = &cobra.Command{
	Use:          "restart",
	Short:        "Restart a happy stack deployment, leaving everything else the same",
	SilenceUsage: true,
	PreRunE: happyCmd.Validate(
		cobra.ExactArgs(1),
		happyCmd.IsStackNameDNSCharset,
		happyCmd.IsStackNameAlphaNumeric,
		func(cmd *cobra.Command, args []string) error {
			checklist := util.NewValidationCheckList()
			return util.ValidateEnvironment(cmd.Context(), checklist.AwsInstalled)
		},
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		happyClient, err := makeHappyClient(cmd, sliceName, stackName, []string{tag}, createTag)
		if err != nil {
			return errors.Wrap(err, "initializing the happy client")
		}

		ctx := context.WithValue(cmd.Context(), options.DryRunKey, dryRun)
		err = validate(
			validateGitTree(happyClient.HappyConfig.GetProjectRoot()),
			validateStackNameAvailable(ctx, happyClient.StackService, stackName, force),
			validateStackExistsUpdate(ctx, stackName, happyClient),
		)
		if err != nil {
			return errors.Wrap(err, "validating happy client")
		}

		backend, err := happyClient.AWSBackend.GetComputeBackend(cmd.Context())
		if err != nil {
			return err
		}

		k8s, ok := backend.(*aws.K8SComputeBackend)
		if !ok {
			return errors.New("not a k8s backend, nothing to do")
		}

		opts := []hapi.APIClientOption{}
		if baseURL != "" {
			opts = append(opts, hapi.WithBaseURL(baseURL))
		}
		api := hapi.MakeAPIClient(happyClient.HappyConfig, happyClient.AWSBackend, opts...)
		configs, err := api.ListConfigs(happyClient.HappyConfig.App(), happyClient.HappyConfig.GetEnv(), stackName)
		if err != nil {
			return errors.Wrap(err, "listing configs")
		}

		var wg sync.WaitGroup
		for _, serviceName := range happyClient.HappyConfig.GetData().Services {
			err := updateDeploymentConfigAndWait(ctx, stackName, serviceName, k8s, &wg, configs)
			if err != nil {
				return errors.Wrap(err, "updating deployment with config")
			}
		}
		wg.Wait()
		logrus.Info("all deployments finished")
		return nil
	},
}

func updateDeploymentConfig(ctx context.Context, stackName, serviceName string, k8s *aws.K8SComputeBackend, configs model.WrappedResolvedAppConfigsWithCount) error {
	deploymentName := k8s.GetDeploymentName(stackName, serviceName)
	deploymentsClient := k8s.ClientSet.AppsV1().Deployments(k8s.KubeConfig.Namespace)
	deployment, err := deploymentsClient.Get(ctx, deploymentName, meta.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "getting deployment %s", deploymentName)
	}
	now := time.Now().Format("20060102150405")
	deployment.ObjectMeta.Annotations["happy/restartedAt"] = now
	deployment.Spec.Template.ObjectMeta.Annotations["happy/restartedAt"] = now
	logrus.Infof("updating deployment %s:%s with happy config", k8s.KubeConfig.Namespace, deploymentName)
	envsKeyIndex := map[string]int{}
	for i := range deployment.Spec.Template.Spec.Containers {
		for i, env := range deployment.Spec.Template.Spec.Containers[i].Env {
			envsKeyIndex[env.Name] = i
		}
	}

	for _, config := range configs.Records {
		// if we are adding a new config value, append it
		_, ok := envsKeyIndex[config.Key]
		if !ok {
			logrus.Debugf("adding new env var %s:%s to deployment %s", config.Key, config.Value, deploymentName)
			for i := range deployment.Spec.Template.Spec.Containers {
				deployment.Spec.Template.Spec.Containers[i].Env = append(deployment.Spec.Template.Spec.Containers[i].Env, core.EnvVar{
					Name:  config.Key,
					Value: config.Value,
				})
			}
			continue
		}

		// otherwise, overwrite the existing value
		for i := range deployment.Spec.Template.Spec.Containers {
			logrus.Debugf("updating env var %s from %s to %s on deployment %s", config.Key, deployment.Spec.Template.Spec.Containers[i].Env[envsKeyIndex[config.Key]].Value, config.Value, deploymentName)
			deployment.Spec.Template.Spec.Containers[i].Env[envsKeyIndex[config.Key]] = core.EnvVar{
				Name:  config.Key,
				Value: config.Value,
			}
		}
	}

	_, err = deploymentsClient.Update(ctx, deployment, meta.UpdateOptions{})
	if err != nil {
		return errors.Wrapf(err, "updating deployment %s env variables", deploymentName)
	}
	return nil
}

func updateDeploymentConfigAndWait(ctx context.Context, stackName, serviceName string, k8s *aws.K8SComputeBackend, wg *sync.WaitGroup, configs model.WrappedResolvedAppConfigsWithCount) error {
	err := updateDeploymentConfig(ctx, stackName, serviceName, k8s, configs)
	if err != nil {
		return err
	}

	waitDeploymentReady(ctx, stackName, serviceName, k8s, wg)
	return nil
}

func waitDeploymentReady(ctx context.Context, stackName, serviceName string, k8s *aws.K8SComputeBackend, wg *sync.WaitGroup) {
	deploymentName := k8s.GetDeploymentName(stackName, serviceName)
	deploymentsClient := k8s.ClientSet.AppsV1().Deployments(k8s.KubeConfig.Namespace)

	wg.Add(1)
	go func() {
		defer wg.Done()
		waiter := time.Tick(5 * time.Second)
		for range waiter {
			d, err := deploymentsClient.Get(ctx, deploymentName, meta.GetOptions{})
			if err != nil {
				logrus.Error(err)
			}
			if d.Status.UnavailableReplicas != 0 {
				logrus.Infof("waiting for deployment %s to be ready (waiting for %d unavailable pods)", deploymentName, d.Status.UnavailableReplicas)
				continue
			}
			break
		}
		logrus.Infof("deployment %s is restarted", deploymentName)
	}()
}
