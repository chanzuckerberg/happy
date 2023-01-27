package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	kube "github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type StacklistK8S struct{}

type StacklistK8SBackend struct {
	clientSet kubernetes.Interface
	config    *rest.Config
	k8sConfig kube.K8SConfig
}

const (
	awsApiCallMaxRetries   = 100
	awsApiCallBackoffDelay = time.Second * 5
)

func MakeStacklistK8SBackend(ctx context.Context, payload model.AppStackPayload2) (*StacklistK8SBackend, error) {
	options := []func(*configv2.LoadOptions) error{
		configv2.WithRegion(payload.AwsRegion),
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

	payload.K8SConfig.AuthMethod = "eks"
	clientSet, config, err := kube.CreateK8sClient(ctx, *payload.K8SConfig, clients, kube.DefaultK8sClientCreator)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create k8s client")
	}

	return &StacklistK8SBackend{
		clientSet: clientSet,
		config:    config,
		k8sConfig: *payload.K8SConfig,
	}, nil
}

func (s *StacklistK8S) GetAppStacks(ctx context.Context, payload model.AppStackPayload2) ([]*model.AppStack, error) {
	fmt.Println("...here in k8s GetAppStacks")
	backend, err := MakeStacklistK8SBackend(ctx, payload)
	if err != nil {
		return nil, err
	}

	fmt.Println("...getting param")
	paramOutput, err := backend.getParam(ctx)
	if err != nil {
		return nil, err
	}

	var stacklist []string
	err = json.Unmarshal([]byte(paramOutput), &stacklist)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json")
	}

	stacks := []*model.AppStack{}
	for _, stackName := range stacklist {
		appStack := model.MakeAppStack(payload.AppName, payload.Environment, stackName)
		stacks = append(stacks, &appStack)
	}

	return stacks, nil
}

func (s *StacklistK8SBackend) getParam(ctx context.Context) (string, error) {
	configMap, err := s.clientSet.CoreV1().ConfigMaps(s.k8sConfig.Namespace).Get(ctx, "stacklist", v1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve stacklist configmap")
	}

	if value, ok := configMap.Data["stacklist"]; ok {
		return value, nil
	}

	return "", errors.Wrapf(err, "unable to retrieve a stacklist key from stacklist configmap")
}
