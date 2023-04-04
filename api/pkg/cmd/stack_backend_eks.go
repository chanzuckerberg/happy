package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	compute_backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/k8s"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type StackBackendEKS struct{}

type EKSBackendClient struct {
	clientSet kubernetes.Interface
	config    *rest.Config
	k8sConfig kube.K8SConfig
}

const (
	awsApiCallMaxRetries   = 100
	awsApiCallBackoffDelay = time.Second * 5
)

func MakeEKSBackendClient(ctx context.Context, payload model.AppStackPayload) (*EKSBackendClient, error) {
	options := []func(*configv2.LoadOptions) error{
		configv2.WithRegion(payload.AwsRegion),
		configv2.WithCredentialsProvider(request.MakeCredentialProvider(ctx)),
		configv2.WithRetryer(func() aws.Retryer {
			// Unless specified, we run into ThrottlingException when repeating calls, when following logs or waiting on a condition.
			return retry.AddWithMaxBackoffDelay(retry.AddWithMaxAttempts(retry.NewStandard(), awsApiCallMaxRetries), awsApiCallBackoffDelay)
		}),
	}

	if payload.AwsProfile != "" {
		options = append(options, configv2.WithSharedConfigProfile(payload.AwsProfile))
	}

	conf, err := configv2.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create an aws session")
	}

	eksclient := eks.NewFromConfig(conf)
	stspresignclient := sts.NewPresignClient(sts.NewFromConfig(conf))

	clients := kube.AwsClients{
		EksClient:        eksclient,
		StsPresignClient: stspresignclient,
	}

	k8sConfig := k8s.K8SConfig{
		Namespace:  payload.K8SNamespace,
		ClusterID:  payload.K8SClusterId,
		AuthMethod: "eks",
	}
	clientSet, config, err := kube.CreateK8sClient(ctx, k8sConfig, clients, kube.DefaultK8sClientCreator)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create k8s client")
	}

	return &EKSBackendClient{
		clientSet: clientSet,
		config:    config,
		k8sConfig: k8sConfig,
	}, nil
}

func (s *StackBackendEKS) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStack, error) {
	backend, err := MakeEKSBackendClient(ctx, payload)
	if err != nil {
		return nil, err
	}

	computeBackend := compute_backend.K8SComputeBackend{
		ClientSet: backend.clientSet,
	}

	paramOutput, err := computeBackend.GetParam(ctx, "stacklist")
	if err != nil {
		return nil, err
	}

	stacks, err := convertParamToStacklist(paramOutput, payload)
	if err != nil {
		return nil, err
	}

	integrationSectet, _, err := computeBackend.GetIntegrationSecret(ctx)

	workspaceRepo := workspace_repo.NewWorkspaceRepo(
		integrationSectet.Tfe.Url,
		integrationSectet.Tfe.Org,
	)

	for _, stack := range stacks {
		workspace, err := workspaceRepo.GetWorkspace(ctx, fmt.Sprintf("%s-%s", payload.AppMetadata.Environment, stack.AppName))
		if err != nil {
			stack.WorkspaceUrl := workspace.GetWorkspaceUrl()
		}
	}

	return stacks, nil
}

// func (s *EKSBackendClient) getParam(ctx context.Context) (string, error) {
// 	computeBackend := compute_backend.K8SComputeBackend{
// 		ClientSet: s.clientSet,
// 	}
// 	return computeBackend.GetParam(ctx, "stacklist")
// }

// func (s *EKSBackendClient) getIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error) {
// 	computeBackend := compute_backend.K8SComputeBackend{
// 		ClientSet: s.clientSet,
// 	}
// 	return computeBackend.GetIntegrationSecret(ctx)
// }
