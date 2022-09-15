package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type K8SComputeBackend struct {
	Backend     *Backend
	ClientSet   *kubernetes.Clientset
	HappyConfig *config.HappyConfig
}

func NewK8SComputeBackend(happyConfig *config.HappyConfig, b *Backend) (interfaces.ComputeBackend, error) {
	clusterId := happyConfig.K8SConfig().ClusterID

	// TOOD: Add the ability to select a cluster by context from .kube/config
	clusterInfo, err := b.eksclient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: &clusterId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to get k8s cluster configuration")
	}
	logrus.Infof("K8S Cluster: %s (%s)\n", *(clusterInfo.Cluster).Name, *(clusterInfo.Cluster).Version)

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
	rawConfig, err := clientcmd.NewDefaultClientConfig(config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create kube config")
	}
	rawConfig.BearerToken = getAuthToken(*b.awsConfig, clusterId)

	clientset, _ := kubernetes.NewForConfig(rawConfig)

	return &K8SComputeBackend{
		Backend:     b,
		ClientSet:   clientset,
		HappyConfig: happyConfig,
	}, nil
}

func getAuthToken(config aws.Config, clusterName string) string {
	stsClient := sts.NewFromConfig(config)
	stsSigner := sts.NewPresignClient(stsClient)
	presignedURLRequest, _ := stsSigner.PresignGetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{}, func(presignOptions *sts.PresignOptions) {
		presignOptions.ClientOptions = append(presignOptions.ClientOptions, func(stsOptions *sts.Options) {
			// Add clusterId Header
			stsOptions.APIOptions = append(stsOptions.APIOptions, smithyhttp.SetHeaderValue(clusterIDHeader, clusterName))
			// Add X-Amz-Expires query param
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
