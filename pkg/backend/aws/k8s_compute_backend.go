package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

type K8SComputeBackend struct {
	Backend     *Backend
	ClientSet   *kubernetes.Clientset
	HappyConfig *config.HappyConfig
}

func NewK8SComputeBackend(ctx context.Context, happyConfig *config.HappyConfig, b *Backend) (interfaces.ComputeBackend, error) {
	var rawConfig *rest.Config
	if happyConfig.K8SConfig().AuthMethod == "eks" {
		// Constructs client configuration dynamically
		clusterId := happyConfig.K8SConfig().ClusterID

		clusterInfo, err := b.eksclient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
			Name: &clusterId,
		})
		if err != nil {
			return nil, errors.Wrap(err, "unable to get k8s cluster configuration")
		}
		logrus.Infof("EKS Authenticated K8S Cluster: %s (%s)\n", *(clusterInfo.Cluster).Name, *(clusterInfo.Cluster).Version)

		cert, _ := base64.RawStdEncoding.DecodeString(*clusterInfo.Cluster.CertificateAuthority.Data)
		config := clientcmdapi.Config{
			APIVersion: "v1",
			Kind:       "Config",
			Clusters: map[string]*clientcmdapi.Cluster{
				"cluster": {
					Server:                   *clusterInfo.Cluster.Endpoint,
					CertificateAuthorityData: cert,
				},
			},
			Contexts: map[string]*clientcmdapi.Context{
				"cluster": {
					Cluster: "cluster",
				},
			},
			CurrentContext: "cluster",
		}
		rawConfig, err = clientcmd.NewDefaultClientConfig(config, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return nil, errors.Wrap(err, "unable to create kubeconfig")
		}
		rawConfig.BearerToken = getAuthToken(ctx, b, clusterId)
	} else if happyConfig.K8SConfig().AuthMethod == "kubeconfig" {
		// Uses a context from kubeconfig file
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		var err error
		rawConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{
				CurrentContext: happyConfig.K8SConfig().Context,
			}).ClientConfig()
		if err != nil {
			return nil, errors.Wrap(err, "unable to detect cluster configuration")
		}
		logrus.Info("Kubeconfig Authenticated K8S Cluster\n")
	} else {
		return nil, errors.New("unsupported authentication type")
	}

	clientset, err := kubernetes.NewForConfig(rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to instantiate k8s client")
	}

	return &K8SComputeBackend{
		Backend:     b,
		ClientSet:   clientset,
		HappyConfig: happyConfig,
	}, nil
}

func getAuthToken(ctx context.Context, b *Backend, clusterName string) string {
	presignedURLRequest, _ := b.stspresignclient.PresignGetCallerIdentity(ctx, &sts.GetCallerIdentityInput{}, func(presignOptions *sts.PresignOptions) {
		presignOptions.ClientOptions = append(presignOptions.ClientOptions, func(stsOptions *sts.Options) {
			stsOptions.APIOptions = append(stsOptions.APIOptions, smithyhttp.SetHeaderValue(clusterIDHeader, clusterName))
			stsOptions.APIOptions = append(stsOptions.APIOptions, smithyhttp.SetHeaderValue("X-Amz-Expires", "90"))
		})
	})
	return v1Prefix + base64.RawURLEncoding.EncodeToString([]byte(presignedURLRequest.URL))
}

func (b *K8SComputeBackend) GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error) {
	secret, err := b.ClientSet.CoreV1().Secrets(b.HappyConfig.K8SConfig().Namespace).Get(ctx, "integration-secret", v1.GetOptions{})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to retrieve integration secret")
	}

	if value, ok := secret.Data["integration_secret"]; ok {
		secret := &config.IntegrationSecret{}
		err = json.Unmarshal(value, secret)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not json parse integraiton secret")
		}
		arn := ""
		return secret, &arn, nil
	}
	return nil, nil, errors.New("integration-secret key is missing from the integration secret")
}
