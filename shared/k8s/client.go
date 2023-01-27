package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8SConfig struct {
	Namespace      string `yaml:"namespace"`
	ClusterID      string `yaml:"cluster_id"`       // used with the 'eks' auth_method
	AuthMethod     string `yaml:"auth_method"`      // 'eks' or 'kubeconfig'; 'eks' will construct auth dynamically, 'kubeconfig' will re-use the context from kube-config file.
	Context        string `yaml:"context"`          // used with the kubeconfig auth_method
	KubeConfigPath string `yaml:"kube_config_path"` // used with the kubeconfig auth_method
}

type K8sClientCreator func(config *rest.Config) (kubernetes.Interface, error)

// func CreateK8sClient(ctx context.Context, k8sConfig K8SConfig) (kubernetes.Interface, error) {
// 	var rawConfig *rest.Config
// 	var err error
// 	if k8sConfig.AuthMethod == "eks" {
// 		// Constructs client configuration dynamically
// 		clusterId := k8sConfig.ClusterID
// 		rawConfig, err = createEKSConfig(ctx, clusterId, b)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "unable to create kubeconfig using EKS cluster id")
// 		}
// 		rawConfig.BearerToken = getAuthToken(ctx, b, clusterId)
// 	} else if k8sConfig.AuthMethod == "kubeconfig" {
// 		// Uses a context from kubeconfig file
// 		kubeconfig := strings.TrimSpace(k8sConfig.KubeConfigPath)
// 		if len(kubeconfig) == 0 {
// 			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
// 		}
// 		rawConfig, err = createK8SConfig(kubeconfig, k8sConfig.Context)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "unable to create kubeconfig using kubernetes context name")
// 		}
// 		logrus.Info("Kubeconfig Authenticated K8S Cluster\n")
// 	} else {
// 		return nil, errors.New("unsupported authentication type")
// 	}

// 	clientset, err := clientCreator(rawConfig)
// }
