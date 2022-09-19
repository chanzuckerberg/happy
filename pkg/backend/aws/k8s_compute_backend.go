package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"strings"

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

type k8sClientCreator func(config *rest.Config) (kubernetes.Interface, error)

type K8SComputeBackend struct {
	Backend     *Backend
	ClientSet   kubernetes.Interface
	HappyConfig *config.HappyConfig
}

func NewK8SComputeBackend(ctx context.Context, happyConfig *config.HappyConfig, b *Backend, clientCreator k8sClientCreator) (interfaces.ComputeBackend, error) {
	var rawConfig *rest.Config
	var err error
	if happyConfig.K8SConfig().AuthMethod == "eks" {
		// Constructs client configuration dynamically
		clusterId := happyConfig.K8SConfig().ClusterID
		rawConfig, err = createEKSConfig(clusterId, b)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create kubeconfig using EKS cluster id")
		}
		rawConfig.BearerToken = getAuthToken(ctx, b, clusterId)
	} else if happyConfig.K8SConfig().AuthMethod == "kubeconfig" {
		// Uses a context from kubeconfig file
		kubeconfig := strings.TrimSpace(happyConfig.K8SConfig().KubeConfigPath)
		if len(kubeconfig) == 0 {
			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}
		rawConfig, err = createK8SConfig(kubeconfig, happyConfig.K8SConfig().Context)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create kubeconfig using kubernetes context name")
		}
		logrus.Info("Kubeconfig Authenticated K8S Cluster\n")
	} else {
		return nil, errors.New("unsupported authentication type")
	}

	clientset, err := clientCreator(rawConfig)
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

func (k8s *K8SComputeBackend) GetIntegrationSecret(ctx context.Context) (*config.IntegrationSecret, *string, error) {
	secret, err := k8s.ClientSet.CoreV1().Secrets(k8s.HappyConfig.K8SConfig().Namespace).Get(ctx, "integration-secret", v1.GetOptions{})
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
	return nil, nil, errors.New("integration_secret key is missing from the integration secret")
}

func createEKSConfig(clusterId string, b *Backend) (*rest.Config, error) {
	var rawConfig *rest.Config
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
	return rawConfig, nil
}

func createK8SConfig(kubeconfig string, context string) (*rest.Config, error) {
	rawConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to detect cluster configuration")
	}
	return rawConfig, nil
}
