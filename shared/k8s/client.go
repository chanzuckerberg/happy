package k8s

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/chanzuckerberg/happy/shared/aws/interfaces"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

const (
	clusterIDHeader = "x-k8s-aws-id"
	v1Prefix        = "k8s-aws-v1."
)

const AuthMethodEKS string = "eks"
const AuthMethodKubeConfig string = "kubeconfig"

type K8SConfig struct {
	Namespace      string `yaml:"namespace" json:"namespace,omitempty"`
	ClusterID      string `yaml:"cluster_id" json:"cluster_id,omitempty"`             // used with the 'eks' auth_method
	AuthMethod     string `yaml:"auth_method" json:"auth_method,omitempty"`           // 'eks' or 'kubeconfig'; 'eks' will construct auth dynamically, 'kubeconfig' will re-use the context from kube-config file.
	Context        string `yaml:"context" json:"context,omitempty"`                   // used with the kubeconfig auth_method
	KubeConfigPath string `yaml:"kube_config_path" json:"kube_config_path,omitempty"` // used with the kubeconfig auth_method
}

type AwsClients struct {
	EksClient        interfaces.EKSAPI
	StsPresignClient interfaces.STSPresignAPI
}

type K8sClientCreator func(config *rest.Config) (kubernetes.Interface, error)

func DefaultK8sClientCreator(config *rest.Config) (kubernetes.Interface, error) {
	return kubernetes.NewForConfig(config)
}

func CreateK8sClient(ctx context.Context, k8sConfig K8SConfig, awsClients AwsClients, clientCreator K8sClientCreator) (kubernetes.Interface, *rest.Config, error) {
	var rawConfig *rest.Config
	var err error
	if k8sConfig.AuthMethod == AuthMethodEKS {
		// Constructs client configuration dynamically
		clusterId := k8sConfig.ClusterID
		rawConfig, err = CreateEKSConfig(ctx, awsClients.EksClient, clusterId)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to create kubeconfig using EKS cluster id")
		}
		rawConfig.BearerToken = GetAuthToken(ctx, awsClients.StsPresignClient, clusterId)
	} else if k8sConfig.AuthMethod == AuthMethodKubeConfig {
		// Uses a context from kubeconfig file
		kubeconfig := strings.TrimSpace(k8sConfig.KubeConfigPath)
		if len(kubeconfig) == 0 {
			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}
		rawConfig, err = CreateK8SConfig(kubeconfig, k8sConfig.Context)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to create kubeconfig using kubernetes context name")
		}
		logrus.Info("Kubeconfig Authenticated K8S Cluster\n")
	} else {
		return nil, nil, errors.Errorf("unsupported authentication type: %s", k8sConfig.AuthMethod)
	}

	clientset, err := clientCreator(rawConfig)
	return clientset, rawConfig, err
}

func CreateEKSConfig(ctx context.Context, eksclient interfaces.EKSAPI, clusterId string) (*rest.Config, error) {
	out, err := eksclient.ListClusters(ctx, &eks.ListClustersInput{})
	if err != nil {
		logrus.Errorf("Unable to list EKS clusters: %s", err.Error())
	} else {
		logrus.Debug("Found clusters:")
		for _, cluster := range out.Clusters {
			logrus.Debugf("  Cluster ID: %s", cluster)
		}
	}
	var rawConfig *rest.Config
	clusterInfo, err := eksclient.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: &clusterId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "eks:DescribeCluster `%s` failed. unable to get k8s cluster configuration", clusterId)
	}
	logrus.Debugf("EKS Authenticated K8S Cluster: %s (%s)\n", *(clusterInfo.Cluster).Name, *(clusterInfo.Cluster).Version)

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

func CreateK8SConfig(kubeconfig string, context string) (*rest.Config, error) {
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

func GetAuthToken(ctx context.Context, stsClient interfaces.STSPresignAPI, clusterName string) string {
	presignedURLRequest, _ := stsClient.PresignGetCallerIdentity(ctx, &sts.GetCallerIdentityInput{}, func(presignOptions *sts.PresignOptions) {
		presignOptions.ClientOptions = append(presignOptions.ClientOptions, func(stsOptions *sts.Options) {
			stsOptions.APIOptions = append(stsOptions.APIOptions, smithyhttp.SetHeaderValue(clusterIDHeader, clusterName))
			stsOptions.APIOptions = append(stsOptions.APIOptions, smithyhttp.SetHeaderValue("X-Amz-Expires", "90"))
		})
	})
	return v1Prefix + base64.RawURLEncoding.EncodeToString([]byte(presignedURLRequest.URL))
}
