package k8s

type K8SConfig struct {
	Namespace      string `yaml:"namespace"`
	ClusterID      string `yaml:"cluster_id"`       // used with the 'eks' auth_method
	AuthMethod     string `yaml:"auth_method"`      // 'eks' or 'kubeconfig'; 'eks' will construct auth dynamically, 'kubeconfig' will re-use the context from kube-config file.
	Context        string `yaml:"context"`          // used with the kubeconfig auth_method
	KubeConfigPath string `yaml:"kube_config_path"` // used with the kubeconfig auth_method
}
