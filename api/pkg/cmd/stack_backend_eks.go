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
	"github.com/chanzuckerberg/happy/shared/k8s"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	paramOutput, err := backend.getParam(ctx)
	if err != nil {
		return nil, err
	}

	return convertParamToStacklist(paramOutput, payload)
}

func (s *EKSBackendClient) getParam(ctx context.Context) (string, error) {
	configMap, err := s.clientSet.CoreV1().ConfigMaps(s.k8sConfig.Namespace).Get(ctx, "stacklist", v1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve stacklist configmap")
	}

	if value, ok := configMap.Data["stacklist"]; ok {
		return value, nil
	}

	return "", errors.Wrapf(err, "unable to retrieve a stacklist key from stacklist configmap")
}
