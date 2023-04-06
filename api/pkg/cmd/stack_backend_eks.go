package cmd

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	compute_backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type StackBackendEKS struct{}

type EKSBackendClient struct {
	clientSet kubernetes.Interface
	config    *rest.Config
	k8sConfig k8s.K8SConfig
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

	clients := k8s.AwsClients{
		EksClient:        eksclient,
		StsPresignClient: stspresignclient,
	}

	k8sConfig := k8s.K8SConfig{
		Namespace:  payload.K8SNamespace,
		ClusterID:  payload.K8SClusterId,
		AuthMethod: "eks",
	}
	clientSet, config, err := k8s.CreateK8sClient(ctx, k8sConfig, clients, k8s.DefaultK8sClientCreator)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create k8s client")
	}

	return &EKSBackendClient{
		clientSet: clientSet,
		config:    config,
		k8sConfig: k8sConfig,
	}, nil
}

func (s *StackBackendEKS) GetAppStacks(ctx context.Context, payload model.AppStackPayload) ([]*model.AppStackResponse, error) {
	backend, err := MakeEKSBackendClient(ctx, payload)
	if err != nil {
		return nil, err
	}

	computeBackend := compute_backend.K8SComputeBackend{
		ClientSet:  backend.clientSet,
		KubeConfig: backend.k8sConfig,
	}

	paramOutput, err := computeBackend.GetParam(ctx, "stacklist")
	if err != nil {
		return nil, err
	}

	integrationSecret, _, err := computeBackend.GetIntegrationSecret(ctx)
	if err != nil {
		return nil, err
	}

	stacklist, err := parseParamToStacklist(paramOutput)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json")
	}

	return enrichStacklistMetadata(ctx, stacklist, payload, integrationSecret)
}
