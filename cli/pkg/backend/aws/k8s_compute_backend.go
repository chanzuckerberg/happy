package aws

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces"
	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	dockerterm "github.com/moby/term"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/util/term"
)

type k8sClientCreator func(config *rest.Config) (kubernetes.Interface, error)

type K8SComputeBackend struct {
	Backend     *Backend
	ClientSet   kubernetes.Interface
	HappyConfig *config.HappyConfig
	rawConfig   *rest.Config
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
		rawConfig:   rawConfig,
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

func (k8s *K8SComputeBackend) GetParam(ctx context.Context, name string) (string, error) {
	configMap, err := k8s.ClientSet.CoreV1().ConfigMaps(k8s.HappyConfig.K8SConfig().Namespace).Get(ctx, "stacklist", v1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve stacklist configmap")
	}

	if value, ok := configMap.Data["stacklist"]; ok {
		return value, nil
	}

	return "", errors.Wrapf(err, "unable to retrieve a stacklist key from stacklist configmap")
}

func (k8s *K8SComputeBackend) WriteParam(
	ctx context.Context,
	name string,
	val string,
) error {
	configMap, err := k8s.ClientSet.CoreV1().ConfigMaps(k8s.HappyConfig.K8SConfig().Namespace).Get(ctx, "stacklist", v1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve stacklist configmap")
	}
	configMap.Data["stacklist"] = val
	_, err = k8s.ClientSet.CoreV1().ConfigMaps(k8s.HappyConfig.K8SConfig().Namespace).Update(ctx, configMap, v1.UpdateOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to update stacklist configmap")
	}
	return nil
}

func (k8s *K8SComputeBackend) PrintLogs(ctx context.Context, stackName string, serviceName string, opts ...util.PrintOption) error {
	deploymentName := fmt.Sprintf("%s-%s", stackName, serviceName)
	labelSelector := v1.LabelSelector{MatchLabels: map[string]string{"app": deploymentName}}
	pods, err := k8s.ClientSet.CoreV1().Pods(k8s.HappyConfig.K8SConfig().Namespace).List(ctx, v1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	})
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve a list of pods for deployment %s", deploymentName)
	}
	logrus.Infof("Found %d matching pods.", len(pods.Items))

	for _, pod := range pods.Items {
		logrus.Infof("... streaming logs from pod %s ...", pod.Name)

		logs, err := k8s.ClientSet.CoreV1().Pods(k8s.HappyConfig.K8SConfig().Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
			Follow: false,
		}).Stream(ctx)
		if err != nil {
			return errors.Wrapf(err, "unable to retrieve logs from pod %s", pod.Name)
		}
		defer logs.Close()
		reader := bufio.NewScanner(logs)

		for reader.Scan() {
			logrus.Info(string(reader.Bytes()))
		}
		logrus.Infof("... done streaming ...")
	}
	return nil
}

func (k8s *K8SComputeBackend) RunTask(ctx context.Context, taskDefArn string, launchType config.LaunchType) error {
	// TODO: not implemented
	return errors.New("not implemented")
}

func (k8s *K8SComputeBackend) Shell(ctx context.Context, stackName string, serviceName string) error {
	deploymentName := fmt.Sprintf("%s-%s", stackName, serviceName)
	labelSelector := v1.LabelSelector{MatchLabels: map[string]string{"app": deploymentName}}
	pods, err := k8s.ClientSet.CoreV1().Pods(k8s.HappyConfig.K8SConfig().Namespace).List(ctx, v1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	})
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve a list of pods for deployment %s", deploymentName)
	}
	if len(pods.Items) == 0 {
		return errors.New("No matching pods found")
	}
	logrus.Infof("Found %d matching pods.", len(pods.Items))

	pod, err := k8s.ClientSet.CoreV1().Pods(k8s.HappyConfig.K8SConfig().Namespace).Get(ctx, pods.Items[0].Name, v1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve pod information for %s", pods.Items[0].Name)
	}

	if len(pod.Spec.Containers) > 1 {
		return errors.New("There's more than one container in a pod")
	}

	containerName := pod.Spec.Containers[0].Name

	req := k8s.ClientSet.CoreV1().RESTClient().Post().Resource("pods").Name(pod.Name).Namespace(pod.Namespace).SubResource("exec").Param("container", containerName)

	eo := &corev1.PodExecOptions{
		Container: containerName,
		Command:   strings.Fields("sh"),
		Stdout:    true,
		Stdin:     true,
		Stderr:    false,
		TTY:       true,
	}

	req.VersionedParams(eo, scheme.ParameterCodec)
	logrus.Info(req.URL())

	exec, err := remotecommand.NewSPDYExecutor(k8s.rawConfig, http.MethodPost, req.URL())
	if err != nil {
		return errors.WithStack(err)
	}

	stdin, stdout, stderr := dockerterm.StdStreams()
	streamOptions := remotecommand.StreamOptions{
		Stdout: stdout,
		Stderr: stderr,
		Tty:    true,
	}
	streamOptions.Stdin = stdin
	t := term.TTY{
		In:  stdin,
		Out: stdout,
		Raw: true,
	}
	streamOptions.TerminalSizeQueue = t.MonitorSize(t.GetSize())
	return t.Safe(func() error { return exec.Stream(streamOptions) })
}
